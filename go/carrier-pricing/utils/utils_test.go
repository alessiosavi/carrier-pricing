package utils

import "testing"

var veichleConf string = `../conf/veichle_data.json`
var carrierConf string = `../conf/carrier_data.json`

func TestInitVehicleList(t *testing.T) {
	vehicle := InitVehicleList(veichleConf)
	t.Log(vehicle)
}

func TestInitCarrierList(t *testing.T) {
	carrier := InitCarrierList(carrierConf)
	t.Log(carrier)
}

func TestAddPercent(t *testing.T) {
	if AddPercent(100, 10) != 110 {
		t.Fail()
	}
}
