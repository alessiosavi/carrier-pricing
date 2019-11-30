package mongoutils

import (
	"carrier-pricing/datastructure"
	"log"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoClient struct {
	M *mgo.Session
}

// InitMongoDBConnection return a session to the Mongo instances configured in input.
// If input is null connect to the default instances
func (m *MongoClient) InitMongoDBConnection(host string, port int, connectionMode string, refreshMode bool) {
	log.Println("Connecting to MongoDB using: ", host)
	var err error
	if stringutils.IsBlank(host) {
		host = "0.0.0.0"
	}

	m.M, err = mgo.Dial(host) // Connection to MongoDB
	if err != nil {
		log.Fatalln("Error! MongoDB does not reply! :/")
	}
	// Verify mongo connection
	err = m.M.Ping()
	if err != nil {
		log.Fatal("Unable to connect to mongo!", err)
	}
	log.Println("Connection init finished succesfully! | ", m.M)
	m.M.SetMode(mgo.Strong, true) // Configuring MongoDB session
}

// RemoveCollectionFromDB is used for remove the collection in input from the db
func (m *MongoClient) RemoveCollectionFromDB(database string, collection string) error {
	log.Println("Removing collection: ", collection, " From DB: ", database)
	err := m.M.DB(database).C(collection).DropCollection()
	if err != nil {
		log.Fatalln("Fatal error :/ ", err)
		return err
	}
	log.Println("Collection ", collection, " removed succesfully!")
	return nil
}

// PopulateData is delegated to initialize the MongoDB collection related to the promotion of the carriers
func (m *MongoClient) PopulateData(carrier []datastructure.CarrierList, database, collection string) {
	for i := range carrier {
		if m.InsertData(carrier[i], database, collection, carrier[i].CarrierName) == "ALREDY_EXIST" {
			log.Println("Carrier [" + carrier[i].CarrierName + "] alredy exists!")
		}
	}
}

// QueryVehicle is deleated to extract from the DB the list of carriers that are compliant with the given veichle
// NOTE: validation on the veichle is made during the HTTP api validation checking on Redis
func (m *MongoClient) QueryVehicle(database, collection, vehicle string) []datastructure.PriceList {
	//db.carriers_name.find({"services.vehicles":{"$in":["large_van"]}})
	t := []string{vehicle}
	var carrier []datastructure.CarrierList
	var priceList []datastructure.PriceList
	// Extract the promotion related to the input vehicle
	err := m.M.DB(database).C(collection).Find(bson.M{"services.vehicles": bson.M{"$in": t}}).All(&carrier)
	if err != nil {
		log.Println("Error during search: [" + err.Error() + "]")
		return priceList
	}

	// Populating price list
	for i := range carrier {
		for j := range carrier[i].Services {
			for k := range carrier[i].Services[j].Vehicles {
				if carrier[i].Services[j].Vehicles[k] == vehicle {
					var service datastructure.PriceList
					service.Service = carrier[i].CarrierName
					service.DeliveryTime = carrier[i].Services[j].DeliveryTime
					//NOTE: Save markup % in price
					service.Price = carrier[i].Services[j].Markup
					priceList = append(priceList, service)
				}
			}
		}
	}

	return priceList
}

// InsertData is delegated to populate the given MongoDB collection using the carrier data in input
func (m *MongoClient) InsertData(carrier datastructure.CarrierList, database, collection, username string) string {
	log.Println("Verify if carrier is alredy registered in DB: ", database, " | Collection: ", collection, " | Carrier: ", carrier)
	err := m.M.DB(database).C(collection).Find(bson.M{"carriername": username}).Select(bson.M{"carriername": 1}).One(nil) // Searching the user

	if err == mgo.ErrNotFound { // User is not present into the DB
		log.Println("Registering new carrier ...")
		err = m.M.DB(database).C(collection).Insert(carrier)
		if err != nil {
			log.Fatalln("Some error occurs during insert, impossible to register a new carrier :/ | Err: ", err)
			return "KO"
		}
		log.Println("carrier registered! -> ", carrier)
		return "OK"
	} else if err != nil {
		log.Fatalln("Error during query. Err:", err.Error())
	}
	log.Println("carrier alredy exists! | ", carrier, " avoiding to overwrite")
	return "ALREDY_EXIST"
}
