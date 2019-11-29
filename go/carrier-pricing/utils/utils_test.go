package utils

import "testing"

var jsonConf string = `../conf/veichle_data.json`

func TestInitVehicleList(t *testing.T) {
	vehicle := InitVehicleList(jsonConf)
	t.Log(vehicle)
}
