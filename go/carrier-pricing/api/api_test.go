package api

import (
	"carrier-pricing/datastructure"
	"encoding/json"
	"testing"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
	requests "github.com/alessiosavi/Requests"
	requeststruct "github.com/alessiosavi/Requests/datastructure"
)

func TestWrongPickupPostcode(t *testing.T) {
	testData := []string{"S11A1AA", "SWQA1AA", "SW1AAAA", "SW1A1A1", "SW1A11A"}
	var url string = `https://localhost:8080/quotes`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	for _, test := range testData {
		body := `{"pickup_postcode":"` + test + `","delivery_postcode":"SW1A1AA"}`
		resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
		if resp.StatusCode == 200 {
			err := json.Unmarshal(resp.Body, &respStruct)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if respStruct.Error != "PICKUP_POSTCODE_INVALID" {
				t.Error(respStruct.Error)
			}
		} else {
			t.Log(resp.Error)
			t.Fail()
		}

	}
}

func TestWrongDeliveryPostcode(t *testing.T) {
	testData := []string{"S11A1AA", "SWQA1AA", "SW1AAAA", "SW1A1A1", "SW1A11A"}
	var url string = `https://localhost:8080/quotes`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	for _, test := range testData {
		body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"` + test + `"}`
		resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
		if resp.StatusCode == 200 {
			err := json.Unmarshal(resp.Body, &respStruct)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if respStruct.Error != "DELIVERY_POSTCODE_INVALID" {
				t.Error(respStruct.Error)
			}
		} else {
			t.Log(resp.Error)
			t.Fail()
		}

	}
}

func TestQuote(t *testing.T) {
	var url string = `https://localhost:8080/quotes`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if !stringutils.IsBlank(respStruct.Error) {
			t.Error("Unexpected error [", respStruct.Error, "]")
		}
		if respStruct.Price != 316 {
			t.Error("Wrong price [", respStruct.Price, "]")
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestVeichleNotManaged(t *testing.T) {
	var url string = `https://localhost:8080/vehicle`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT","vehicle":"a"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "VEHICLE_NOT_MANAGED" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestVeichleNotProvided(t *testing.T) {
	var url string = `https://localhost:8080/vehicle`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "VEHICLE_NOT_PROVIDED" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestDeliveryPostEmpty(t *testing.T) {
	var url string = `https://localhost:8080/quotes`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "PICKUP_POSTCODE_INVALID" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestPostCodeEmpty(t *testing.T) {
	var url string = `https://localhost:8080/quotes`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "PICKUP_POSTCODE_INVALID" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestNoVeichle(t *testing.T) {
	var url string = `https://localhost:8080/vehicle`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "VEHICLE_NOT_PROVIDED" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestRequestNotPost(t *testing.T) {
	var url string = `https://localhost:8080/vehicle`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT"}`
	resp = requests.SendRequest(url, "GET", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if respStruct.Error != "REQ_NOT_POST" {
			t.Error(respStruct.Error)
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}

func TestVeichleOK(t *testing.T) {
	var url string = `https://localhost:8080/vehicle`
	var resp *requeststruct.RequestResponse
	var respStruct datastructure.ResponseQuotes

	body := `{"pickup_postcode":"SW1A1AA","delivery_postcode":"EC2A3LT","vehicle":"bicycle"}`
	resp = requests.SendRequest(url, "POST", nil, []byte(body), true)
	if resp.StatusCode == 200 {
		err := json.Unmarshal(resp.Body, &respStruct)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if !stringutils.IsBlank(respStruct.Error) {
			t.Error("Unexpected error!")
		}
		if respStruct.Price != 348 {
			t.Error("Wrong price [", respStruct.Price, "]")
		}

	} else {
		t.Log(resp.Error)
		t.Fail()
	}
}
