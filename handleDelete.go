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

func handleDelete(client gocloak.GoCloak, event data.LocationEvent, ctx context.Context, token gocloak.JWT, realm string, reader kafka.Reader, m kafka.Message) {
	fmt.Println("received new delete message")
	var entity data.Entity
	err := json.Unmarshal(event.Entity, &entity) //TODO: Error handling
	if err != nil {
		fmt.Println(err)
		return
	}
	entity.Name = s.ToLower(entity.Name)

	searchLocation, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{})
	if len(searchLocation) != 0 {
		hasDeleted := false
		for _, subgroup := range *searchLocation[0].SubGroups {
			locationIds := (*subgroup.Attributes)["locationId"]
			if len(locationIds) != 0 && locationIds[0] == strconv.Itoa(entity.Id) {
				err = client.DeleteGroup(ctx, token.AccessToken, realm, *subgroup.ID) //TODO: Error handling
				hasDeleted = true
				if err != nil {
					fmt.Println(err)
					break
				}
			}

		}
		if hasDeleted == false {
			fmt.Println("tried to delete location that doesn't exist")

		}
		reader.CommitMessages(ctx, m)
	}
}
