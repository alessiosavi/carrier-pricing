package datastructure

// RequestQuotes is delegated to marshal the request related to the `/quotes` endpoint
type RequestQuotes struct {
	PickupPostcode   string `json:"pickup_postcode"`
	DeliveryPostcode string `json:"delivery_postcode"`
	Veichle          string `json:"vehicle"`
}

// RequestQuotes is delegated to marshal the response related to the `/quotes` endpoint
type ResponseQuotes struct {
	PickupPostcode   string      `json:"pickup_postcode"`
	DeliveryPostcode string      `json:"delivery_postcode"`
	Price            int         `json:"price"`
	PriceList        []PriceList `json:"price_list,omitempty"`
	Error            string      `json:"error,omitempty"`
}

// Configuration save the data necessary for initialize the tool
type Configuration struct {
	Host    string `json:"Host"`
	Port    int    `json:"Port"`
	Version string `json:"Version"`
	SSL     struct {
		Path    string `json:"Path"`
		Cert    string `json:"Cert"`
		Key     string `json:"Key"`
		Enabled bool   `json:"Enabled"`
	} `json:"SSL"`
	Redis struct {
		Host string `json:"Host"`
		Port string `json:"Port"`
		DB   int    `json:"DB"`
	} `json:"Redis"`
	Mongo struct {
		Host    string `json:"Host"`
		Port    int    `json:"Port"`
		Carrier struct {
			DB         string `json:"DB"`
			Collection string `json:"Collection"`
		} `json:"Carrier"`
	} `json:"Mongo"`
}

// VehicleList is delegated to save the % of price increment
type VehicleList struct {
	Vehicles []string `json:"vehicles"`
	Prices   []int    `json:"prices"`
}

type CarrierList struct {
	CarrierName string `json:"carrier_name"`
	BasePrice   int    `json:"base_price"`
	Services    []struct {
		DeliveryTime int      `json:"delivery_time"`
		Markup       int      `json:"markup"`
		Vehicles     []string `json:"vehicles"`
	} `json:"services"`
}

type PriceList struct {
	Service      string `json:"service"`
	Price        int    `json:"price"`
	DeliveryTime int    `json:"delivery_time"`
}
