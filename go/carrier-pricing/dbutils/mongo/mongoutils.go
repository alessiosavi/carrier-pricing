package mongoutils

import (
	"carrier-pricing/datastructure"
	"log"

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
	m.M, err = mgo.Dial(host) // Connection to MongoDB
	if err != nil {
		log.Fatalln("Error! MongoDB does not reply! :/", m.M)
	}
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

func (m *MongoClient) PopulateData(carrier []datastructure.CarrierList, database, collection string) {
	for i := range carrier {
		m.InsertData(carrier[i], database, collection, carrier[i].CarrierName)
	}
}

func (m *MongoClient) QueryVehicle(database, collection, vehicle string) []datastructure.PriceList {
	//db.carriers_name.find({"services.vehicles":{"$in":["large_van"]}})
	t := []string{vehicle}
	var carrier []datastructure.CarrierList
	err := m.M.DB(database).C(collection).Find(bson.M{"services.vehicles": bson.M{"$in": t}}).All(&carrier)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Carrier -> ", carrier)
	var priceList []datastructure.PriceList
	for i := range carrier {
		for j := range carrier[i].Services {
			for k := range carrier[i].Services[j].Vehicles {
				if carrier[i].Services[j].Vehicles[k] == "small_van" {
					var service datastructure.PriceList
					service.Service = carrier[i].CarrierName
					service.DeliveryTime = carrier[i].Services[j].DeliveryTime
					// Save markup in price
					service.Price = carrier[i].Services[j].Markup
					priceList = append(priceList, service)
				}
			}
		}
	}

	return priceList
}

// InsertData is used for insert a generic data into a collection
// It take in input the session, database and collection where insert the change
func (m *MongoClient) InsertData(carrier datastructure.CarrierList, database, collection, username string) string {
	log.Println("Verify if carrier is alredy registered in DB: ", database, " | Collection: ", collection, " | Carrier: ", carrier)
	err := m.M.DB(database).C(collection).Find(bson.M{"carriername": username}).Select(bson.M{"carriername": 1}).One(nil) // Searching the user

	log.Println("Error ->", err)
	if err != nil { // User is not present into the DB
		log.Println("Registering new carrier ...")
		err = m.M.DB(database).C(collection).Insert(carrier)
		if err != nil {
			log.Fatalln("Some error occurs during insert, impossible to register a new carrier :/ | Err: ", err)
			return "KO"
		}
		log.Println("carrier registered! -> ", carrier)
		return "OK"
	}
	log.Println("carrier alredy exists! | ", carrier, " avoiding to overwrite")
	return "ALREDY_EXIST"
}

// RemoveUser Remove a registered user from MongoDB
func (m *MongoClient) RemoveUser(database, collection, carrier_name string) error {
	log.Println("RemoveUser | Removing user: ", carrier_name, " | From DB: ", database, " | Collection: ", collection)
	err := m.M.DB(database).C(collection).Remove(bson.M{"Username": carrier_name})
	if err != nil {
		log.Fatalln("RemoveUser | Error during delete of user :( | User: ", carrier_name, " | Session: ", m.M, " | Error: ", err)
		return err
	}
	log.Println("RemoveUser | Correctly deleted | ", carrier_name)
	return nil
}
