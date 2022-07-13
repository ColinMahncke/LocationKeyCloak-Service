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

func handleUpdate(client gocloak.GoCloak, event data.LocationEvent, ctx context.Context, token gocloak.JWT, realm string, reader kafka.Reader, m kafka.Message) {
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
		hasUpdated := false
		for _, subgroup := range *searchLocation[0].SubGroups {
			locationIds := (*subgroup.Attributes)["locationId"]
			if len(locationIds) != 0 && locationIds[0] == strconv.Itoa(entity.Id) {
				*subgroup.Name = entity.Name
				err = client.UpdateGroup(ctx, token.AccessToken, realm, subgroup)

				if err != nil {
					fmt.Println(err)
					break
				}

				for _, usergroup := range *subgroup.SubGroups {
					if s.HasPrefix(*usergroup.Name, "admin_") {
						*usergroup.Name = "admin_" + entity.Name
						err = client.UpdateGroup(ctx, token.AccessToken, realm, usergroup)

						if err != nil {
							fmt.Println(err)
							break
						}
					}
					if s.HasPrefix(*usergroup.Name, "api_") {
						*usergroup.Name = "api_" + entity.Name
						err = client.UpdateGroup(ctx, token.AccessToken, realm, usergroup)

						if err != nil {
							fmt.Println(err)
							break
						}
					}

				}
				hasUpdated = true

			}

		}
		if hasUpdated == false {
			fmt.Println("tried to update location that doesn't exist")

		}
		reader.CommitMessages(ctx, m)
	}
}
