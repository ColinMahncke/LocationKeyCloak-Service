package main

import (
	"context"
	"encoding/json"
	"fmt"
	"locationKeycloakService/data"
	"strconv"

	s "strings"

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
		json.Unmarshal(m.Value, &event)
		if event.EventType == "Create" {
			fmt.Println("received new create message")
			var entity data.Entity
			json.Unmarshal(event.Entity, &entity)
			entity.Name = s.ToLower(entity.Name)
			locationPath := entity.Name
			searchLocation, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{Search: &locationPath})
			if len(searchLocation) != 0 {
				reader.CommitMessages(ctx, m)
				fmt.Println("tried to create location that already exist")
				continue
			}

			searchParent, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{Search: &parentName})

			if len(searchParent) != 0 {

				attributes := make(map[string][]string)
				attributes["locationId"] = []string{strconv.FormatInt(int64(entity.Id), 10)}
				newID, err := client.CreateChildGroup(ctx, token.AccessToken, realm, *searchParent[0].ID, gocloak.Group{Name: &entity.Name, Attributes: &attributes})
				if err != nil {
					fmt.Println(err)

				}
				apiName := "api_" + entity.Name
				_, _ = client.CreateChildGroup(ctx, token.AccessToken, realm, newID, gocloak.Group{Name: &apiName, Attributes: &attributes})
				if err != nil {
					fmt.Println(err)
				}
				adminName := "admin_" + entity.Name
				_, _ = client.CreateChildGroup(ctx, token.AccessToken, realm, newID, gocloak.Group{Name: &adminName, Attributes: &attributes})
				if err != nil {
					fmt.Println(err)

				}
				reader.CommitMessages(ctx, m)
			}
		}
		if event.EventType == "Delete" {
			fmt.Println("received a new message")
			var entity data.Entity
			json.Unmarshal(event.Entity, &entity)
			entity.Name = s.ToLower(entity.Name)
			locationPath := entity.Name

			searchLocation, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{Search: &locationPath})
			if len(searchLocation) != 0 {
				client.DeleteGroup(ctx, token.AccessToken, realm, *searchLocation[0].ID)
			}
			reader.CommitMessages(ctx, m)
		}
	}

}
