package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"aws_test/common"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/websocket"
)

// Store recent status retrieved from sensors
var currentStatus struct {
	InternalTemp float64
	ExternalTemp float64
	LightStatus bool
}

var wsConnections []*websocket.Conn

// MQTT stuff
var mqttHandler *common.MQTTHandler

// Handles all received messages
func onMQTTMessageReceived(topic string, payload []byte) {
	value := string(payload[:])

	switch topic {
	case "sensors/temp/external":
		currentStatus.ExternalTemp, _ = strconv.ParseFloat(value, 64)
	case "sensors/temp/internal":
		currentStatus.InternalTemp, _ = strconv.ParseFloat(value, 64)
	case "sensors/light/internal":
		currentStatus.LightStatus, _ = strconv.ParseBool(value)
	default:
		fmt.Println("Unknown MQTT message received:")
		fmt.Printf("Topic: %s\n", topic)
		fmt.Printf("Message: %s\n", payload)
	}

	// TODO: update UI immediately
	for i,v := range wsConnections {
		if err := v.WriteJSON("{}"); err != nil {
			fmt.Println(err)
			// TODO: Remove failed connection from update list
			wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
		}
	}
}

// Holds authentication data
func createMQTTHandler() *common.MQTTHandler {
	// Initialize MQTT connection
	return common.NewMQTTHandler(
		"ssl://a1x6d6ym1e2e7r.iot.eu-central-1.amazonaws.com:8883",
		"GatewayThing",
		"gateway/certs/AWS_Root_CA.pem",
		"gateway/certs/gateway.pem.crt",
		"gateway/certs/gateway.private.pem.key",
		[]string{ "sensors/temp/+", "sensors/light/+"},
		onMQTTMessageReceived,
	)
}

// REST API handler
func FileHandler(w http.ResponseWriter, r *http.Request) {
	// Return UI part
	vars := mux.Vars(r)
	fileName := vars["fileName"]
	if fileName == "" {
		fileName = "index.html"
	}

	html, err := ioutil.ReadFile("gateway/html/" + fileName)
	if err != nil {
		http.NotFound(w, r)
	}
	w.Write(html)
}

// REST API handler
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Return recent status
	json.NewEncoder(w).Encode(currentStatus)
}

// REST API handler
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Set new update interval by posting MQTT message to sensors
	vars := mux.Vars(r)
	interval := vars["interval"]
	mqttHandler.Publish("sensors/setting/interval", interval)
}

// REST API handler
func LightStateHandler(w http.ResponseWriter, r *http.Request) {
	// Update light enabled state by posting MQTT message to sensors
	vars := mux.Vars(r)
	newState := vars["state"]
	mqttHandler.Publish("sensors/setting/light", newState)
}

// Main GO entry point
func main() {
	mqttHandler = createMQTTHandler()

	router := mux.NewRouter()
	router.HandleFunc("/", FileHandler)
	router.HandleFunc("/{fileName}", FileHandler)
	router.HandleFunc("/api/status", StatusHandler)
	router.HandleFunc("/api/interval/{interval:[0-9]+}", UpdateHandler)
	router.HandleFunc("/api/light/{state:(?:on|off)}", LightStateHandler)

	// TODO: Websocket test
	router.HandleFunc("/ws/one", wsHandler)

	fmt.Println("Gateway initialized")
	fmt.Println(http.ListenAndServe(":8080", router))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	/*
	if r.Header.Get("Origin") != "http://" + r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	*/

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	wsConnections = append(wsConnections, conn)
}