package utils

import "testing"

var jsonConf string = `../conf/veichle_data.json`

func TestInitVehicleList(t *testing.T) {
	vehicle := InitVehicleList(jsonConf)
	t.Log(vehicle)
}

func TestAddPercent(t *testing.T) {
	if AddPercent(100, 10) != 110 {
		t.Fail()
	}
}
