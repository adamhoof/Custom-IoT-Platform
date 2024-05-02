package mqtt_handlers

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func OnConnectionLost(client MQTT.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
