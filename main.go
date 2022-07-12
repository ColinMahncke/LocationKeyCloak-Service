package main

import (
	"context"
	"encoding/json"
	"fmt"
	"locationKeycloakService/data"

	"github.com/Nerzal/gocloak"
	"github.com/segmentio/kafka-go"
)

func main() {
	Load()
	host := Getenv("host", "localhost:9093")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{host},
		Topic:     "location_events",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer reader.Close()
	realm := Getenv("realmName", "realmName")
	client := gocloak.NewClient(Getenv("keycloakHost", "https://mycool.keycloak.instance"))
	token, err := client.LoginAdmin(Getenv("user", "user"), Getenv("password", "password"), realm)
	if err != nil {
		fmt.Println(err)
		return
	}

	parentGroup, err := client.GetGroup(token.AccessToken, realm, "leap_on_the_beach")
	if err != nil {
		fmt.Println(err)
		return
	}
	if parentGroup == nil {
		err = client.CreateGroup(token.AccessToken, realm, gocloak.Group{Name: "leap_on_the_beach"})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		var event data.LocationEvent
		json.Unmarshal(m.Value, &event)
		if event.EventType == "Create" {
			var entity data.CreateEntity
			json.Unmarshal(event.Entity, &entity)

		}
	}

}
