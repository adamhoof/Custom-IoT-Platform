package main

import (
	"NSI-semester-work/internal/mqtt_handlers"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

func setupMqttClient() *MQTT.Client {
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
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe(os.Getenv("MQTT_LOGIN_REQUEST_TOPIC"), 1, mqtt_handlers.HandleLogin)
	client.Subscribe(os.Getenv("MQTT_POST_TOPIC"), 1, mqtt_handlers.HandlePost)

	return &client
}

func main() {
	mqttClient := setupMqttClient()
	(*mqttClient).Disconnect(250)
}
