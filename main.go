package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Variables set via flag
var (
	// debug is if we should be printing out debug logs
	debug bool
	// host is the host to listen on
	host string
	// port is the port to listen on
	port string
)

type collections struct {
	minimal *mongo.Collection
	full    *mongo.Collection
}

type postData struct {
	Level   string       `json:"level,omitempty"`
	Minimal identityMin  `json:"minimal,omitempty"`
	Full    identityFull `json:"full,omitempty"`
}

type identityMin struct {
	Platform          interface{} `json:"platform,omitempty"`
	OriginalOSVersion interface{} `json:"original_os_version,omitempty"`
	CurrentOSVersion  interface{} `json:"current_os_version,omitempty"`
	InstanceType      interface{} `json:"instance_type,omitempty"`
}

type identityFull struct {
	Hardware    interface{} `json:"hardware,omitempty"`
	Network     interface{} `json:"network,omitempty"`
	ContainerRT interface{} `json:"container_rt,omitempty"`
}

type responseData struct {
	Success bool `json:"success"`
}

var collectionsPtr *collections = nil

// connect to MongoDB
func connectMongoDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// disconnect from MongoDB
func disconnectMongoDB(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}

// delete the specified db
func deleteDBFromMongoDB(db *mongo.Database) {
	log.Println("Dropped database: ", db.Name())
	db.Drop(context.TODO())
}

// write the response and sets the status
func writeResponse(w *http.ResponseWriter, res responseData, status int) {
	resJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(*w, "Failed to parse struct `responseData` into JSON object", http.StatusInternalServerError)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(status)
	(*w).Write(resJSON)
}

// logMiddleware produces simple logs when wrapped around a HandlerFunc
func logMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s - %s", r.RemoteAddr, r.Method, r.URL.Path)
		if debug {
			log.Printf("Request: %+v", r)
		}
		// Execute the original handler
		handler(w, r)
	}
}

// main handler for data collected by pinger
func dataHandler(w http.ResponseWriter, r *http.Request) {

	// connect to MongoDB when running test
	if collectionsPtr == nil {
		client := connectMongoDB()
		log.Println("Connected to MongoDB!")
		defer disconnectMongoDB(client)
		defer deleteDBFromMongoDB(client.Database("fcos_pinger_test"))

		// Get handles for `minimal` and `full` collections
		minimalCollection := client.Database("fcos_pinger_test").Collection("minimal")
		fullCollection := client.Database("fcos_pinger_test").Collection("full")
		collections := collections{minimal: minimalCollection, full: fullCollection}
		collectionsPtr = &collections
	}

	// enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// note: to enable CORS, it should handle OPTIONS as well
	if r.Method == "OPTIONS" {
		return
	}

	// If it's a POST, process
	if r.Method == "POST" {
		pd := postData{}
		res := responseData{false}

		contentType := r.Header.Get("Content-type")
		if !strings.Contains(contentType, "application/json") {
			writeResponse(&w, res, http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&pd)
		if err != nil {
			log.Println(err)
			writeResponse(&w, res, http.StatusInternalServerError)
			return
		}

		// Checks if POST json contains a `level` or `minimal` key
		emptyMin := identityMin{nil, nil, nil, nil}
		if pd.Level == "" || pd.Minimal == emptyMin {
			log.Printf("failed: `level` and `minmal` must be set")
			writeResponse(&w, res, http.StatusBadRequest)
			return
		}

		// Checks if request body has reasonal length
		maxLen := 31415
		maxLenStr := os.Getenv("PINGER_MAX_LENGTH")
		if len(maxLenStr) > 0 {
			maxLen, _ = strconv.Atoi(maxLenStr)
		}
		if len(fmt.Sprintf("%v", pd.Level)) > maxLen ||
			len(fmt.Sprintf("%v", pd.Minimal)) > maxLen ||
			len(fmt.Sprintf("%v", pd.Full)) > maxLen {
			log.Printf("failed: request body `level`, `minimal` or `full` too long")
			writeResponse(&w, res, http.StatusBadRequest)
			return
		}
		// Process POST request body
		jsonString, err := json.Marshal(pd)
		if err != nil {
			log.Println(err)
			writeResponse(&w, res, http.StatusInternalServerError)
			return
		}
		res.Success = true
		log.Printf(string(jsonString))

		insertResult, err := collectionsPtr.minimal.InsertOne(context.TODO(), pd.Minimal)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted a single document in `minimal`: ", insertResult.InsertedID)

		insertResult, err = collectionsPtr.full.InsertOne(context.TODO(), pd.Full)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted a single document in `full`: ", insertResult.InsertedID)

		writeResponse(&w, res, http.StatusOK)
		return
	}

	// Otherwise method is not allowed
	log.Println("Method not allowed")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

// main is the main entry point on the CLI
func main() {

	client := connectMongoDB()
	log.Println("Connected to MongoDB!")
	defer disconnectMongoDB(client)

	// Get handles for `minimal` and `full` collections
	minimalCollection := client.Database("fcos_pinger").Collection("minimal")
	fullCollection := client.Database("fcos_pinger").Collection("full")
	collections := collections{minimal: minimalCollection, full: fullCollection}
	collectionsPtr = &collections

	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.StringVar(&host, "host", "127.0.0.1", "Host to listen on")
	flag.StringVar(&port, "port", "5000", "Port to listen on")
	flag.Parse()

	// create the host:port string for use
	listenAddress := fmt.Sprintf("%s:%s", host, port)
	if debug {
		log.Printf("Listening on %s", listenAddress)
	}

	// Map `/` to our dataHandler and wrap it in the log middleware
	http.Handle("/", logMiddleware(http.HandlerFunc(dataHandler)))

	// Run forever on all interfaces on port 5000
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
