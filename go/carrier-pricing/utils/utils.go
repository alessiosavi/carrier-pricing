package utils

import (
	"carrier-pricing/datastructure"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

var regexpFormat string = `[A-Z]{2}[0-9][A-Z][0-9][A-Z]{2}`

// Base36 is delegated to compute the base36 of the givens string and round to the lower integer
// Note: the sanity check are performed during the input validation by the HTTP Api
func Base36(s1, s2 string) int {
	a, _ := strconv.ParseInt(s1, 36, 64)
	b, _ := strconv.ParseInt(s2, 36, 64)
	// Negative value not allowed
	r1 := math.Abs(float64(a - b))
	// Round in defect
	r1 = (r1 / 100000000) - 0.5
	r := math.RoundToEven(r1)
	// We are not an ONG, set minimum price in case of same postal code :D
	if r < 4 {
		r = 4
	}
	// Safe due to high division
	return int(r)
}

func InitRegexp() *regexp.Regexp {
	// Initialize regexp instead of compile every time
	reg, err := regexp.Compile(regexpFormat)
	if err != nil {
		log.Fatalln("Unable to initialize regexp [" + err.Error() + "]")
	}
	return reg
}

func ValidatePostCode(code string, reg *regexp.Regexp) bool {
	return reg.MatchString(code)
}

var configDefaultLocation string = `./conf/conf.json`
var vehicleDefaultLocation string = `./conf/veichle_data.json`
var carrierDefaultLocation string = `./conf/carrier_data.json`

// VerifyCommandLineInput is delegated to manage the inputer parameter provide with the input flag from command line
func VerifyCommandLineInput() (datastructure.Configuration, datastructure.VehicleList, []datastructure.CarrierList) {
	conf := flag.String("config", configDefaultLocation, "Specify the configuration file.")
	vehicle := flag.String("vehicle", vehicleDefaultLocation, "Specify the vehicle:price list.")
	carrier := flag.String("carrier", carrierDefaultLocation, "Specify the carriers json list.")
	flag.Parse()

	if strings.Compare(*conf, configDefaultLocation) == 0 {
		log.Println("Running with default conf [" + *conf + "]. Use `--config file.json` to overwrite the configuration")
	}
	if strings.Compare(*vehicle, vehicleDefaultLocation) == 0 {
		log.Println("Running with default conf [" + *vehicle + "]. Use `--vehicle file.json` to overwrite the configuration")
	}
	if strings.Compare(*vehicle, carrierDefaultLocation) == 0 {
		log.Println("Running with default conf [" + *carrier + "]. Use `--carrier file.json` to overwrite the configuration")
	}

	if !fileutils.FileExists(*conf) {
		log.Fatalln("File [" + *conf + "] does not exists!")
	}

	// Reading global conf
	file, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatalln("VerifyCommandLineInput | can't open config file: ", err)
	}
	cfg := datastructure.Configuration{}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalln("VerifyCommandLineInput | can't decode config JSON in ["+*conf+"]: ", err)
	}
	log.Println("VerifyCommandLineInput | Conf loaded -> ", cfg)

	// Init veichle markup (for Redis)
	vehicles := InitVehicleList(*vehicle)

	// Init Carriers promotion (for MongoDB)
	carriers := InitCarrierList(*carrier)
	return cfg, vehicles, carriers

}

func InitCarrierList(carrierJSON string) []datastructure.CarrierList {

	if !fileutils.FileExists(carrierJSON) {
		log.Fatalln("carrier file not found [" + carrierJSON + "]")
	}

	carriersData, err := ioutil.ReadFile(carrierJSON)
	if err != nil {
		log.Fatal("Unable to read [" + carrierJSON + "] Json")
	}

	var carriers []datastructure.CarrierList
	err = json.Unmarshal(carriersData, &carriers)

	if err != nil {
		log.Fatal("Unable to cast [" + carrierJSON + "] into struct")
	}

	if len(carriers) == 0 {
		log.Fatal("vehicles empty")
	}
	return carriers
}

func InitVehicleList(vehicleJSON string) datastructure.VehicleList {
	if !fileutils.FileExists(vehicleJSON) {
		log.Fatalln("vehicle file not found [" + vehicleJSON + "]")
	}

	vehicleData, err := ioutil.ReadFile(vehicleJSON)
	if err != nil {
		log.Fatal("Unable to read [" + vehicleJSON + "] Json")
	}

	var vehicle datastructure.VehicleList
	err = json.Unmarshal(vehicleData, &vehicle)

	if err != nil {
		log.Fatal("Unable to cast [" + vehicleJSON + "] into struct")
	}

	if len(vehicle.Vehicles) == 0 {
		log.Fatal("vehicles empty")
	}

	return vehicle
}

// ValidateRequestBasic is a method that verify the common data between requests that share the same parameter
func ValidateRequestBasic(sBody string) string {
	// Empty request
	if stringutils.IsBlank(sBody) || sBody == "{}" {
		var e string = "EMPTY_REQUEST"
		return e
	}

	if !strings.Contains(sBody, "delivery_postcode") {
		var e string = "DELIVERY_POSTCODE_EMPTY"
		return e
	}

	if !strings.Contains(sBody, "pickup_postcode") {
		var e string = "PICKUP_POSTCODE_EMPTY"
		return e
	}
	return ""
}

// ValidateRequestBasic is delegated to specify the validation for the PostCode
func ValidatePostCodeRequest(req datastructure.RequestQuotes, reg *regexp.Regexp) string {
	// Manage not valid PickupPostcode
	if !ValidatePostCode(req.PickupPostcode, reg) {
		var e string = "PICKUP_POSTCODE_INVALID"
		return e
	}

	// Manage not valid DeliveryPostcode
	if !ValidatePostCode(req.DeliveryPostcode, reg) {
		var e string = "DELIVERY_POSTCODE_INVALID"
		return e
	}
	return ""
}

// AddPercent is delegated to calculate the increment based on markup
func AddPercent(price int, percent int) int {
	if percent > 0 {
		var p float64 = float64(price*percent) / 100
		total := math.Round(float64(price) + p)
		return int(total)
	}
	return price
}
