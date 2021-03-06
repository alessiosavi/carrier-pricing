package api

import (
	"carrier-pricing/datastructure"
	mongoutils "carrier-pricing/dbutils/mongo"
	redisutils "carrier-pricing/dbutils/redis"
	"carrier-pricing/utils"
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
	helper "github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/valyala/fasthttp"
)

func InitHandler(reg *regexp.Regexp, redis redisutils.RedisClient, mongo mongoutils.MongoClient, db, collection string, certs ...string) (fasthttp.RequestHandler, bool) {

	m := func(ctx *fasthttp.RequestCtx) { // Hook to the API methods "magilogically"
		ctx.Response.Header.Set("carrier-pricing", "v0.0.1") // Set an header just for track the version of the software
		log.Println("REQUEST -->", ctx, "| Headers:", ctx.Request.Header.String())
		tmpChar := "============================================================"
		if ctx.IsPost() {
			log.Println(tmpChar)
			switch string(ctx.Path()) {

			case "/quotes":
				// Allow only POST req
				quotes(ctx, reg)

			case "/vehicle":
				vehicle(ctx, reg, redis)

			case "/carrier":
				carrier(ctx, reg, redis, mongo, db, collection)

			default:
				ctx.SetStatusCode(404)
			}
			log.Println(tmpChar)
		} else {
			var e string = "REQ_NOT_POST"
			manageError(ctx, e)
		}

	}

	var enableSSL bool = true

	if len(certs) != 2 {
		log.Println("Certs not provided, disabling ssl")
		enableSSL = false
	} else {
		// NOTE: Orders of certificate matters
		if !fileutils.FileExists(certs[0]) {
			log.Println("Certificate not provided")
			enableSSL = false
		} else if !fileutils.FileExists(certs[1]) {
			log.Println("KeyFile not provided")
			enableSSL = false
		}
	}

	// The gzipHandler will serve a compress request only if the client request it with headers (Content-Type: gzip, deflate)
	gzipHandler := fasthttp.CompressHandlerLevel(m, fasthttp.CompressBestCompression) // Compress data before sending (if requested by the client)
	return gzipHandler, enableSSL
}
func InitAPIFasthttp(hostname, port string, reg *regexp.Regexp, redisClient redisutils.RedisClient, mongo mongoutils.MongoClient, db, collection string, certs ...string) {

	// Create an handler for the HTTP API service
	gzipHandler, enableSSL := InitHandler(reg, redisClient, mongo, db, collection, certs...)

	s := &fasthttp.Server{
		Handler: gzipHandler,
		// Every response will contain 'Server: carrier-pricing challenge' header.
		Name: "carrier-pricing challenge",
		// set a maxing request size for avoid DDOS
		MaxRequestBodySize: 100 * 1024,
	}

	log.Println("Max size allowed (per file) ->", helper.ConvertSize(float64(s.MaxRequestBodySize), "KB"), "KB")
	// Check if SSL can be enabled
	if enableSSL {
		err := s.ListenAndServeTLS(hostname+":"+port, certs[0], certs[1]) // Try to start the server with input "host:port" received in input
		if err != nil {                                                   // No luck, connection not successfully. Probably port used ...
			if strings.Contains(err.Error(), `PEM inputs may have been switched`) {
				log.Println("WARNING! PEM inputs may have been switched, change the order of certificate files")
			}
			log.Fatalln("Unable to spawn SSL server. Err: " + err.Error())
		}
	} else {
		err := s.ListenAndServe(hostname + ":" + port) // Try to start the server with input "host:port" received in input
		if err != nil {                                // No luck, connection not successfully. Probably port used ...
			log.Fatalln("Unable to spawn server. Err: " + err.Error())
		}
	}
}

func carrier(ctx *fasthttp.RequestCtx, reg *regexp.Regexp, redis redisutils.RedisClient, mongo mongoutils.MongoClient, db, collection string) {

	body := ctx.PostBody()
	e, req := validateVeichleRequest(body, reg, redis)
	if e != "" {
		manageError(ctx, e)
		return
	}

	ok, perc := redis.GetValueFromDB(req.Veichle)
	if !ok {
		var e string = `DB_CONNECTION_ERROR`
		manageError(ctx, e)
		return
	}

	percent, err := strconv.Atoi(perc)
	if err != nil {
		log.Fatal("There is an error in the json file related to the markup increment!")
	}
	// Calculating price
	price := utils.Base36(req.PickupPostcode, req.DeliveryPostcode)
	price = utils.AddPercent(price, percent)

	// Retrieve the list of carriers related to the veichle
	prices := mongo.QueryVehicle(db, collection, req.Veichle)

	if len(prices) == 0 {
		log.Println("No promotion found for veichle [" + req.Veichle + "]")
		var e string = "NO_PROMOTION_FOUND"
		manageError(ctx, e)
		return
	}

	// Calculating additional markup
	for i := range prices {
		prices[i].Price = utils.AddPercent(price, prices[i].Price)
	}

	log.Println("Before sort -> ", prices)
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})
	log.Println("After sort -> ", prices)

	// Populating datastructure
	var resp datastructure.ResponseQuotes
	resp.PickupPostcode = req.PickupPostcode
	resp.DeliveryPostcode = req.DeliveryPostcode
	resp.PriceList = prices
	resp.Price = price

	// Encoding response to stdout
	err = json.NewEncoder(ctx).Encode(resp)
	if err != nil {
		log.Println("Unable to write into customer response")
	}
}
func quotes(ctx *fasthttp.RequestCtx, reg *regexp.Regexp) {
	body := ctx.PostBody()

	e, req := validateQuoteRequest(body, reg)
	if e != "" {
		manageError(ctx, e)
		return
	}

	// Calculating price
	price := utils.Base36(req.PickupPostcode, req.DeliveryPostcode)

	// Populating datastructure
	var resp datastructure.ResponseQuotes
	resp.PickupPostcode = req.PickupPostcode
	resp.DeliveryPostcode = req.DeliveryPostcode
	resp.Price = price

	// Encoding response to stdout
	err := json.NewEncoder(ctx).Encode(resp)
	if err != nil {
		log.Println("Unable to write into customer response")
	}
}

func vehicle(ctx *fasthttp.RequestCtx, reg *regexp.Regexp, redis redisutils.RedisClient) {
	body := ctx.PostBody()
	e, req := validateVeichleRequest(body, reg, redis)
	if e != "" {
		manageError(ctx, e)
		return
	}

	ok, perc := redis.GetValueFromDB(req.Veichle)
	if !ok {
		var e string = `DB_CONNECTION_ERROR`
		manageError(ctx, e)
		return
	}

	percent, _ := strconv.Atoi(perc)
	// Calculating price
	price := utils.Base36(req.PickupPostcode, req.DeliveryPostcode)
	price = utils.AddPercent(price, percent)
	// Populating datastructure
	var resp datastructure.ResponseQuotes
	resp.PickupPostcode = req.PickupPostcode
	resp.DeliveryPostcode = req.DeliveryPostcode
	resp.Price = price

	// Encoding response to stdout
	err := json.NewEncoder(ctx).Encode(resp)
	if err != nil {
		log.Println("Unable to write into customer response")
	}

}

func validateVeichleRequest(body []byte, reg *regexp.Regexp, redis redisutils.RedisClient) (string, datastructure.RequestQuotes) {
	sBody := string(body)
	log.Println("Managing input request: " + sBody)
	// Basic validation
	e := utils.ValidateRequestBasic(sBody)
	if e != "" {
		return e, datastructure.RequestQuotes{}
	}

	if !strings.Contains(sBody, "vehicle") {
		var e string = "VEHICLE_NOT_PROVIDED"
		return e, datastructure.RequestQuotes{}
	}

	var req datastructure.RequestQuotes
	// Cast response into struct
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Unable to parse request [" + sBody + "]. Err: " + err.Error())
		return err.Error(), datastructure.RequestQuotes{}
	}
	log.Println("Request unmarshalled [", req, "] | Validating post codes")
	e = utils.ValidatePostCodeRequest(req, reg)
	if e != "" {
		return e, datastructure.RequestQuotes{}
	}

	if stringutils.IsBlank(req.Veichle) {
		var e string = "VEHICLE_PARM_EMPTY"
		return e, datastructure.RequestQuotes{}
	}

	// Be sure that the vehicle is managed
	if !redis.VerifyKeyExistence(req.Veichle) {
		var e string = "VEHICLE_NOT_MANAGED"
		return e, datastructure.RequestQuotes{}
	}

	return "", req
}

func validateQuoteRequest(body []byte, reg *regexp.Regexp) (string, datastructure.RequestQuotes) {
	sBody := string(body)
	log.Println("Managing input request: " + sBody)
	// Basic validation
	e := utils.ValidateRequestBasic(sBody)
	if e != "" {
		return e, datastructure.RequestQuotes{}
	}

	var req datastructure.RequestQuotes
	// Cast response into struct
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Unable to parse request [" + sBody + "]. Err: " + err.Error())
		return err.Error(), datastructure.RequestQuotes{}
	}
	log.Println("Request unmarshalled [", req, "]")

	e = utils.ValidatePostCodeRequest(req, reg)
	if e != "" {
		return e, datastructure.RequestQuotes{}
	}

	return "", req
}

// manageError is delegated to print the error in the customer response
func manageError(ctx *fasthttp.RequestCtx, e string) {
	log.Println("Printing error [" + e + "] into customer response")
	var resp datastructure.ResponseQuotes
	resp.Error = e
	err := json.NewEncoder(ctx).Encode(resp)
	if err != nil {
		log.Println("Unable to write into customer response")
	}
}
