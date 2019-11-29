package datastructure

// RequestQuotes is delegated to marshal the request related to the `/quotes` endpoint
type RequestQuotes struct {
	PickupPostcode   string `json:"pickup_postcode"`
	DeliveryPostcode string `json:"delivery_postcode"`
	Veichle          string `json:"vehicle"`
}

// RequestQuotes is delegated to marshal the response related to the `/quotes` endpoint
type ResponseQuotes struct {
	PickupPostcode   string `json:"pickup_postcode"`
	DeliveryPostcode string `json:"delivery_postcode"`
	Price            int    `json:"price"`
	Error            string `json:"error,omitempty"`
}

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
}

// VehicleList is delegated to save the % of price increment
type VehicleList struct {
	Vehicles []string `json:"vehicles"`
	Prices   []int    `json:"prices"`
}
