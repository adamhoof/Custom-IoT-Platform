package mqtt_handlers

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func OnReconnectingHandler(client MQTT.Client, opts *MQTT.ClientOptions) {
	log.Println("reconnecting mqtt client")
}
