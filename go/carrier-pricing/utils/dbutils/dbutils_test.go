package dbutils

import (
	"carrier-pricing/utils"
	"testing"
)

var jsonConf string = `../../conf/veichle_data.json`

func TestInitVeichleDB(t *testing.T) {
	var r RedisClient
	r.ConnectToDb("", "", 1)
	vehicle := utils.InitVehicleList(jsonConf)
	r.InitVeichleDB(vehicle)
	r.Shutdown()
}

func TestRemoveVeichleDB(t *testing.T) {
	var r RedisClient
	r.ConnectToDb("", "", 1)
	vehicle := utils.InitVehicleList(jsonConf)
	r.RemoveVeichleDB(vehicle)
	r.Shutdown()
}
