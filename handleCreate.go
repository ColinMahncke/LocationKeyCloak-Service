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

func HandleCreate(client gocloak.GoCloak, event data.LocationEvent, ctx context.Context, token gocloak.JWT, realm string, reader kafka.Reader, m kafka.Message) {
	fmt.Println("received new create message")
	var entity data.Entity
	err := json.Unmarshal(event.Entity, &entity) //TODO: Error handling
	if err != nil {
		fmt.Println(err)
		return
	}
	entity.Name = s.ToLower(entity.Name)
	locationPath := entity.Name
	searchLocation, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{Search: &locationPath})
	if len(searchLocation) != 0 {
		hasCreated := false
		for _, subgroup := range *searchLocation[0].SubGroups {
			if *subgroup.Name == entity.Name {
				hasCreated = true
				break
			}
		}
		if hasCreated == true {
			reader.CommitMessages(ctx, m)
			fmt.Println("tried to create location that already exist")
			return
		}
	}
	parentName := "leap_on_the_beach"
	searchParent, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{Search: &parentName})

	if len(searchParent) != 0 {

		attributes := make(map[string][]string)
		attributes["locationId"] = []string{strconv.FormatInt(int64(entity.Id), 10)}
		newID, err := client.CreateChildGroup(ctx, token.AccessToken, realm, *searchParent[0].ID, gocloak.Group{Name: &entity.Name, Attributes: &attributes})
		if err != nil {
			fmt.Println(err)
			return
		}
		apiName := "api_" + entity.Name
		_, err = client.CreateChildGroup(ctx, token.AccessToken, realm, newID, gocloak.Group{Name: &apiName, Attributes: &attributes})
		if err != nil {
			fmt.Println(err)
			return
		}
		adminName := "admin_" + entity.Name
		_, err = client.CreateChildGroup(ctx, token.AccessToken, realm, newID, gocloak.Group{Name: &adminName, Attributes: &attributes})
		if err != nil {
			fmt.Println(err)
			return
		}
		reader.CommitMessages(ctx, m)
	}
}
