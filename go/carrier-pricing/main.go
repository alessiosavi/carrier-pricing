package main

import (
	"carrier-pricing/api"
	"carrier-pricing/utils"
	"carrier-pricing/utils/dbutils"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	cfg, vehicle := utils.VerifyCommandLineInput()
	log.Println("CFG: ", cfg, "Veichle: ", vehicle)

	// Avoid to initialize regexp for every request
	reg := utils.InitRegexp()

	// Initialize Redis connection
	var r dbutils.RedisClient
	r.ConnectToDb(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	r.InitVeichleDB(vehicle)

	api.InitAPIFasthttp("localhost", "8080", reg, r, `./conf/ssl/localhost.crt`, `./conf/ssl/localhost.key`)

}
