package redisutils

import (
	"carrier-pricing/datastructure"
	"log"
	"strconv"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	R *redis.Client
}

// ConnectToDb use emtpy string for hardcoded port
func (r *RedisClient) ConnectToDb(addr string, port string, db int) {
	if stringutils.IsBlank(addr) {
		addr = "localhost"
	}
	if stringutils.IsBlank(port) {
		port = "6379"
	}
	if db < 0 {
		db = 1
	}

	r.R = redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: "", // no password set
		DB:       db,
	})
	log.Println("Connecting to -> ", r.R)
	err := r.R.Ping().Err()
	if err != nil {
		log.Fatalln("Impossibile to connecto to DB ....| CLIENT: ", addr, ":", port, " | ERR: ", err)
	}
}

// RemoveValueFromDB is delegated to check if a key is alredy inserted and return the value
func (r *RedisClient) RemoveValueFromDB(key string) bool {
	err := r.R.Del(key).Err()
	if err == nil {
		log.Println("RemoveValueFromDB | SUCCESS | Key: ", key, " | Removed")
		return true
	} else if err == redis.Nil {
		log.Println("RemoveValueFromDB | Key -> ", key, " does not exist")
		return false
	}
	log.Println("RemoveValueFromDB | Fatal exception during retrieving of data [", key, "] | Redis: ", r.R)
	return false
}

// InsertIntoClient set the two value into the Databased pointed from the client
func (r *RedisClient) InsertIntoClient(key string, value string) bool {
	log.Println("InsertIntoClient | Inserting -> (", key, ":", value, ")")
	err := r.R.Set(key, value, 0).Err() // Inserting the values into the DB
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("InsertIntoClient | INSERTED SUCCESFULLY!! | (", key, ":", value, ")")
	return true
}

// GetValueFromDB is delegated to check if a key is alredy inserted and return the value
func (r *RedisClient) GetValueFromDB(key string) (bool, string) {
	tmp, err := r.R.Get(key).Result()
	if err == nil {
		log.Println("SUCCESS | Key: ", key, " | Value: ", tmp)
		return true, tmp
	} else if err == redis.Nil {
		log.Println("Key -> ", key, " does not exist")
		return false, ""
	}
	log.Println("GetValueFromDB | Fatal exception during retrieving of data [", key, "] | Redis: ", r.R)
	return false, ""
}

// VerifyKeyExistence is delegated to check if a key is alredy inserted
func (r *RedisClient) VerifyKeyExistence(key string) bool {
	_, err := r.R.Get(key).Result()
	if err == nil {
		return true
	} else if err == redis.Nil {
		log.Println("Key -> ", key, " does not exist")
		return false
	}
	log.Println("GetValueFromDB | Fatal exception during retrieving of data [", key, "] | Redis: ", r.R)
	return false
}

func (r *RedisClient) InitVeichleDB(vehicle datastructure.VehicleList) {
	for i, key := range vehicle.Vehicles {
		if !r.VerifyKeyExistence(key) {
			log.Println("Inserting ["+key+":", vehicle.Prices[i], "]")
			r.InsertIntoClient(key, strconv.Itoa(vehicle.Prices[i]))
		}
	}
}

func (r *RedisClient) RemoveVeichleDB(vehicle datastructure.VehicleList) {
	for _, key := range vehicle.Vehicles {
		if r.VerifyKeyExistence(key) {
			log.Println("Removing [" + key + "]")
			if r.RemoveValueFromDB(key) {
				log.Println("Key [" + key + "] removed!")
			}
		}
	}
}

func (r *RedisClient) Shutdown() {
	r.R.Close()
}
