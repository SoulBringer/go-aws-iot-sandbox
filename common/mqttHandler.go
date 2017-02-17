package common

import (
	"fmt"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	client MQTT.Client
	mqttMessageHandler MQTTMessageHandler
}

type MQTTMessageHandler func (topic string, payload []byte)

// Setup new MQTTHandler instance
func NewMQTTHandler(brokerUri string, clientId string, rootCertificate string, clientCertificate string,
	clientCertificateKey string, subscribeTopics []string, messageHandler MQTTMessageHandler) *MQTTHandler {

	self := &MQTTHandler{}

	// Create a ClientOptions
	clientOptions := MQTT.NewClientOptions()
	clientOptions.AddBroker(brokerUri)
	clientOptions.SetClientID(clientId)
	clientOptions.SetDefaultPublishHandler(self.onMessageReceived)
	clientOptions.SetTLSConfig(createTLSConfig(rootCertificate, clientCertificate, clientCertificateKey))

	self.client = MQTT.NewClient(clientOptions)
	self.mqttMessageHandler = messageHandler

	// Connect to broker
	if token := self.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe and request messages to be delivered
	for _, v := range subscribeTopics {
		if token := self.client.Subscribe(v, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	return self
}

// MQTT message handler
func (self *MQTTHandler)onMessageReceived(client MQTT.Client, msg MQTT.Message) {
	// Pass data to associated handler
	if self.mqttMessageHandler != nil {
		self.mqttMessageHandler(msg.Topic(), msg.Payload())
	}
}

// Send MQTT message
func (self *MQTTHandler)Publish(topic string, payload string) {
	// Publish message and wait for the receipt from the server after sending
	token := self.client.Publish(topic, 0, false, payload)
	token.Wait()
}

// Create transport level config
func createTLSConfig(rootCertificate string, clientCertificate string, clientCertificateKey string) *tls.Config {
	// Setup root CA certificate
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(rootCertificate)
	if err != nil {
		panic(err)
	}
	certPool.AppendCertsFromPEM(caCert)

	// Setup client certificate provided by AWS
	cert, err := tls.LoadX509KeyPair(clientCertificate, clientCertificateKey)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certPool,
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
}