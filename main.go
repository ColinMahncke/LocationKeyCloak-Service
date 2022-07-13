package main

import (
	"context"
	"encoding/json"
	"fmt"
	"locationKeycloakService/data"

	"github.com/Nerzal/gocloak/v11"
	"github.com/segmentio/kafka-go"
)

func main() {
	//loading of environment variables and host
	Load()

	host := Getenv("host", "localhost:9093")
	//initialization of the reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{host},
		Topic:    "location_events",
		GroupID:  "location_service",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	realm := Getenv("realmName", "realmName")
	client := gocloak.NewClient(Getenv("keycloakHost", "https://mycool.keycloak.instance"))
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, Getenv("user", "user"), Getenv("password", "password"), realm)
	if err != nil {
		fmt.Println(err)
		return
	}

	parentName := "leap_on_the_beach"
	_, _ = client.CreateGroup(ctx, token.AccessToken, realm, gocloak.Group{Name: &parentName})

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		//unmarshal and creation of Groups in Keycloak
		var event data.LocationEvent
		err = json.Unmarshal(m.Value, &event) //TODO: Error handling
		if err != nil {
			fmt.Println(err)
			continue
		}
		if event.EventType == "Create" {
			HandleCreate(client, event, ctx, *token, realm, *reader, m)
		}
		if event.EventType == "Delete" {
			handleDelete(client, event, ctx, *token, realm, *reader, m)

		}

	}
}
