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
	// We are not an ONG, set minimum price in case of same postal code
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

// VerifyCommandLineInput is delegated to manage the inputer parameter provide with the input flag from command line
func VerifyCommandLineInput() datastructure.Configuration {
	c := flag.String("config", "./conf/test.json", "Specify the configuration file.")
	flag.Parse()
	if strings.Compare(*c, "./conf/test.json") == 0 {
		log.Println("Running with default conf [" + *c + "]. Use `--config conf/config.json` to overwrite the configuration")
	}
	if !fileutils.FileExists(*c) {
		log.Fatalln("File [" + *c + "] does not exists!")
	}

	file, err := ioutil.ReadFile(*c)
	if err != nil {
		log.Fatalln("VerifyCommandLineInput | can't open config file: ", err)
	}
	cfg := datastructure.Configuration{}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalln("VerifyCommandLineInput | can't decode config JSON in ["+*c+"]: ", err)
	}
	log.Println("VerifyCommandLineInput | Conf loaded -> ", cfg)

	return cfg
}

func InitVehicleList(vehicleJSON string) datastructure.VehicleList {
	if !fileutils.FileExists(vehicleJSON) {
		log.Fatalln("Veichle file not found [" + vehicleJSON + "]")
	}

	veichleData, err := ioutil.ReadFile(vehicleJSON)
	if err != nil {
		log.Fatal("Unable to read [" + vehicleJSON + "] Json")
	}

	var vehicle datastructure.VehicleList
	err = json.Unmarshal(veichleData, &vehicle)

	if err != nil {
		log.Fatal("Unable to cast [" + vehicleJSON + "] into struct")
	}

	if len(vehicle.Vehicles) != len(vehicle.Prices) {
		log.Fatal("Veichles and price have not the same size")
	}

	return vehicle
}
