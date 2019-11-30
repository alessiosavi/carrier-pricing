package main

import (
	"carrier-pricing/api"
	mongoutils "carrier-pricing/dbutils/mongo"
	redisutils "carrier-pricing/dbutils/redis"
	"carrier-pricing/utils"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// f, err := os.OpenFile("log/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Println("Unable to create log file")
	// }
	// defer f.Close()
	// log.SetOutput(f)
	// log.Println("Logging initialized")
	cfg, vehicle, carriers := utils.VerifyCommandLineInput()
	log.Println("CFG: ", cfg, "Veichle: ", vehicle)

	// Init mongo connection
	var mongoClient mongoutils.MongoClient
	mongoClient.InitMongoDBConnection(cfg.Mongo.Host, cfg.Mongo.Port, "", true)
	defer mongoClient.M.Close()
	mongoClient.PopulateData(carriers, cfg.Mongo.Carrier.DB, cfg.Mongo.Carrier.Collection)
	priceList := mongoClient.QueryVehicle(cfg.Mongo.Carrier.DB, cfg.Mongo.Carrier.Collection, "small_van")

	log.Println("services -> ", priceList)

	// Initialize Redis connection
	var r redisutils.RedisClient
	r.ConnectToDb(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	defer r.R.Close()
	r.InitVeichleDB(vehicle)

	// Avoid to initialize regexp for every request
	reg := utils.InitRegexp()

	api.InitAPIFasthttp("localhost", "8080", reg, r, mongoClient, cfg.Mongo.Carrier.DB, cfg.Mongo.Carrier.Collection, `./conf/ssl/localhost.crt`, `./conf/ssl/localhost.key`)

}
