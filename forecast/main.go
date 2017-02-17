package main

import (
	"fmt"
	"aws_test/common"
	"time"
	"strconv"
)

// MQTT stuff
var updateInterval time.Duration = 30 * time.Second
var publishTopic string = "sensors/temp/external"
var subscribeTopic string = "sensors/setting/interval"

// Handles all received messages
func onMQTTMessageReceived(topic string, payload []byte) {
	if topic == subscribeTopic {
		value, _ := strconv.ParseInt(string(payload[:]), 10,64)
		updateInterval = time.Duration(value) * time.Millisecond
		fmt.Printf("New update interval set: %v\n", updateInterval)
	} else {
		fmt.Println("Unknown MQTT message received:")
		fmt.Printf("Topic: %s\n", topic)
		fmt.Printf("Message: %s\n", payload)
	}
}

// Holds authentication data
func createMQTTHandler() *common.MQTTHandler {
	// Initialize MQTT connection
	return common.NewMQTTHandler(
		"ssl://a1x6d6ym1e2e7r.iot.eu-central-1.amazonaws.com:8883",
		"ForecastThing",
		"forecast/certs/AWS_Root_CA.pem",
		"forecast/certs/forecast.pem.crt",
		"forecast/certs/forecast.private.pem.key",
		[]string{subscribeTopic},
		onMQTTMessageReceived,
	)
}

// Main GO entry point
func main() {
	mqttHandler := createMQTTHandler()
	for {
		temp, err := GetExternalTemperature()
		if err == nil {
			mqttHandler.Publish(publishTopic, fmt.Sprintf("%.2f", temp))
			fmt.Println("External temperature data is sent")
		} else {
			fmt.Printf("Error occurred while retrieving external temperature: %s\n", err)
		}
		time.Sleep(updateInterval)
	}
}