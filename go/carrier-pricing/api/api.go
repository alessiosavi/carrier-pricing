package api

import (
	"carrier-pricing/datastructure"
	"carrier-pricing/utils"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
	helper "github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

func InitHandler(reg *regexp.Regexp, certs ...string) (fasthttp.RequestHandler, bool) {

	m := func(ctx *fasthttp.RequestCtx) { // Hook to the API methods "magilogically"
		ctx.Response.Header.Set("carrier-pricing", "v0.0.1") // Set an header just for track the version of the software
		log.Println("REQUEST -->", ctx, "| Headers:", ctx.Request.Header.String())
		tmpChar := "============================================================"
		switch string(ctx.Path()) {
		case "/":
			ctx.SetStatusCode(404)
			log.Println(tmpChar)
		case "/quotes":
			// Allow only POST req
			if ctx.IsPost() {
				quotes(ctx, reg)
			} else {
				var e string = "REQ_NOT_POST"
				manageError(ctx, e)
			}
		}
	}

	var enableSSL bool = true

	if len(certs) != 2 {
		log.Println("Certs not provided, disabling ssl")
		enableSSL = false
	} else {
		// TODO: Orders of certificate matters
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
func InitAPIFasthttp(hostname, port string, reg *regexp.Regexp, redisClient *redis.Client, certs ...string) {

	gzipHandler, enableSSL := InitHandler(reg, certs...)

	s := &fasthttp.Server{
		Handler: gzipHandler,
		// Every response will contain 'Server: carrier-pricing challenge' header.
		Name: "carrier-pricing challenge",
		// Other Server settings may be set here.
		MaxRequestBodySize: 100 * 1024,
	}

	log.Println("Max size allowed (per file) ->", helper.ConvertSize(float64(s.MaxRequestBodySize), "KB"), "KB")
	if enableSSL {
		err := s.ListenAndServeTLS(hostname+":"+port, certs[0], certs[1]) // Try to start the server with input "host:port" received in input
		if err != nil {                                                   // No luck, connection not successfully. Probably port used ...
			if strings.Contains(err.Error(), `PEM inputs may have been switched`) {
				log.Println("WARNING! PEM inputs may have been switched")
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

func quotes(ctx *fasthttp.RequestCtx, reg *regexp.Regexp) {
	var req datastructure.RequestQuotes

	body := ctx.PostBody()
	sBody := string(body)
	log.Println("Managing input request: " + sBody)
	if stringutils.IsBlank(sBody) || sBody == "{}" {
		var e string = "EMPTY_REQUEST"
		manageError(ctx, e)
		return
	}

	if !strings.Contains(sBody, "delivery_postcode") {
		var e string = "Delivery post code is empty"
		manageError(ctx, e)
		return
	}

	if !strings.Contains(sBody, "pickup_postcode") {
		var e string = "Pickup post code is empty"
		manageError(ctx, e)
		return
	}

	// Cast response into struct
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Unable to parse request [" + sBody + "]. Err: " + err.Error())
		manageError(ctx, err.Error())
		return
	}
	log.Println("Request unmarshalled [", req, "]")

	// Manage not valid PickupPostcode
	if !utils.ValidatePostCode(req.PickupPostcode, reg) {
		var e string = "Pickup post code not valid! [" + req.PickupPostcode + "]"
		manageError(ctx, e)
		return
	}

	// Manage not valid DeliveryPostcode
	if !utils.ValidatePostCode(req.DeliveryPostcode, reg) {
		var e string = "Delivery post code not valid! [" + req.DeliveryPostcode + "]"
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
	err = json.NewEncoder(ctx).Encode(resp)
	if err != nil {
		log.Println("Unable to write into customer response")
	}
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
