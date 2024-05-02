package mqtt_handlers

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func OnConnectHandler(client MQTT.Client) {
	optsReader := client.OptionsReader()
	log.Printf("Connected to mqtt broker: %s", optsReader.Servers()[0])
}
