package main

import (
	"NSI-semester-work/internal/http_handlers"
	"NSI-semester-work/internal/mqtt_handlers"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/http"
	"os"
	"time"
)

func setupMqttClient() (MQTT.Client, error) {
	var broker = os.Getenv("MQTT_BROKER")
	var port = os.Getenv("MQTT_PORT")
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%s", broker, port))
	opts.SetClientID("go_web_server_mqtt_client")
	opts.SetOnConnectHandler(mqtt_handlers.OnConnectHandler)
	opts.SetConnectionLostHandler(mqtt_handlers.OnConnectionLost)
	opts.SetReconnectingHandler(mqtt_handlers.OnReconnectingHandler)
	opts.SetAutoReconnect(true)
	opts.SetOrderMatters(false)
	opts.SetMaxReconnectInterval(time.Second * 10)

	client := MQTT.NewClient(opts)
	fmt.Println("connecting to mqtt broker")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("unable to connect to client %s\n", token.Error())
	}

	if token := client.Subscribe(os.Getenv("MQTT_LOGIN_REQUEST_TOPIC"), 1, mqtt_handlers.HandleLogin); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to subscribe to login topic: %v", token.Error())
	}
	if token := client.Subscribe(os.Getenv("MQTT_POST_TOPIC"), 1, mqtt_handlers.HandlePost); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to subscribe to post topic: %v", token.Error())
	}

	return client, nil
}

func setupHttpServer() error {
	serverHostname := os.Getenv("HTTP_SERVER_HOST")
	port := os.Getenv("HTTP_SERVER_PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/", http_handlers.HandleHome)

	fmt.Printf("starting HTTP server: http://%s:%s\n", serverHostname, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", serverHostname, port), mux); err != nil {
		return fmt.Errorf("unable to start server %s\n", err)
	}
	return nil
}

func main() {
	_, err := setupMqttClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = setupHttpServer(); err != nil {
		log.Fatal(err)
	}
}
