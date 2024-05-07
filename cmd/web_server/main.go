package main

import (
	"NSI-semester-work/internal/db"
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

	return client, nil
}

func setupMqttSubscriptionHandlers(client MQTT.Client, database *db.Database) error {
	if token := client.Subscribe("login/request/+", 0, func(client MQTT.Client, msg MQTT.Message) {
		mqtt_handlers.HandleDeviceLogin(client, msg, database)
	}); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to login topic: %v", token.Error())
	}
	if token := client.Subscribe("provide_value/+", 1, func(client MQTT.Client, msg MQTT.Message) { mqtt_handlers.ValueProvidedHandler(msg, database) }); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to post topic: %v", token.Error())
	}
	return nil
}

func setupHttpServer(database *db.Database, mqttClient MQTT.Client) error {
	serverHostname := os.Getenv("HTTP_SERVER_HOST")
	port := os.Getenv("HTTP_SERVER_PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http_handlers.HomeHandler(w, database) })
	mux.HandleFunc("/dashboard_creator", func(w http.ResponseWriter, r *http.Request) { http_handlers.DashboardCreatorHandler(w, database) })
	mux.HandleFunc("/device_features/{id}", func(w http.ResponseWriter, r *http.Request) { http_handlers.DeviceFeaturesHandler(w, r, database) })
	mux.HandleFunc("/create_dashboard", func(w http.ResponseWriter, r *http.Request) { http_handlers.CreateDashboardHandler(w, r, database) })
	mux.HandleFunc("/dashboard/{id}", func(w http.ResponseWriter, r *http.Request) { http_handlers.DisplayDashboardHandler(w, r, database) })
	mux.HandleFunc("/device/{device_id}/command/{action_name}", func(w http.ResponseWriter, r *http.Request) {
		http_handlers.SendCommandHandler(w, r, mqttClient, database)
	})
	mux.HandleFunc("/device/{device_id}/provide_value/{action_name}", func(w http.ResponseWriter, r *http.Request) { http_handlers.GetLastSensorValueHandler(w, r, database) })

	fmt.Printf("starting HTTP server: http://%s:%s\n", serverHostname, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", serverHostname, port), mux); err != nil {
		return fmt.Errorf("unable to start server %s\n", err)
	}
	return nil
}

func main() {
	mqttClient, err := setupMqttClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	dataSourceName := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_DB"),
	)
	database, err := db.NewDatabase(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(database *db.Database) {
		err = database.Close()
		if err != nil {

		}
	}(database)

	err = setupMqttSubscriptionHandlers(mqttClient, database)
	if err != nil {
		log.Fatal(err)
	}

	if err = setupHttpServer(database, mqttClient); err != nil {
		log.Fatal(err)
	}
}
