package main

/*
import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

// MQTT message handler
var onMessageReceived MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// MQTT client assembler
func createMqttClient() MQTT.Client {
	// Setup transport level config
	certpool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("forecast/certs/AWS_Root_CA.pem")
	if err != nil {
		panic(err)
	}
	certpool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("forecast/certs/forecast.pem.crt", "forecast/certs/forecast.private.pem.key")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
		//MinVersion: tls.VersionTLS12,
		//MaxVersion: tls.VersionTLS12,
	}

	// Create a ClientOptions
	clientOptions := MQTT.NewClientOptions()
	clientOptions.AddBroker("ssl://a1x6d6ym1e2e7r.iot.eu-central-1.amazonaws.com:8883")
	clientOptions.SetClientID("ForecastThing")
	clientOptions.SetDefaultPublishHandler(onMessageReceived)
	clientOptions.SetTLSConfig(tlsConfig)

	return MQTT.NewClient(clientOptions)
}

// Reflect openweathermap forecast JSON structure
type ExternalWeather struct {
	Main struct{
		Temp float64
	}
}

// Retreives external temperature based on openweathermap forecast
func getExternalTemperature() (string, error) {
	forecastUrl := "http://api.openweathermap.org/data/2.5/weather?appid=78ac22a76a45e175dbb87e0fb0b38bd6&units=metric&q=Vinnitsya,ua"
	resp, err := http.Get(forecastUrl)
	temp := ""
	if err == nil {
		defer resp.Body.Close()
		//body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			weather := ExternalWeather{}
			err := json.NewDecoder(resp.Body).Decode(&weather)
			if err == nil {
				temp = fmt.Sprintf("%.2f", weather.Main.Temp)
				fmt.Println(temp)
			} else {
				fmt.Println(err)
			}
		}
	}
	return temp, err
}

// Send forecast MQTT message
func sendForecast(mqttClient MQTT.Client) {
	for {
		fmt.Println("Sending forecast message")
		temperature, err := getExternalTemperature()
		if err == nil {
			// Publish message and wait for the receipt from the server after sending
			token := mqttClient.Publish("sensors/forecast", 0, false, temperature)
			token.Wait()
		} else {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}

// Main GO entry point
func main() {
	// Create and start a client
	mqttClient := createMqttClient()
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe and request messages to be delivered
	if token := mqttClient.Subscribe("command/forecast", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	defer func() {
		// Unsubscribe & disconnect
		if token := mqttClient.Unsubscribe("command/forecast"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		mqttClient.Disconnect(250)
		fmt.Println("Disconnected")
	}()

	// TODO: wait till the end of time, add custom logic
	go sendForecast(mqttClient)
	running := make(chan struct{})
	<- running
	//os.Exit(1)
}
*/