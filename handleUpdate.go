package main

import (
	"encoding/json"
	"fmt"
	"locationKeycloakService/data"
	"strconv"
	s "strings"

	"github.com/Nerzal/gocloak/v11"
	"github.com/segmentio/kafka-go"
)

func HandleUpdate(event data.LocationEvent, reader kafka.Reader, m kafka.Message) {

	client, ctx, token, realm, err := GetKeycloakClient()
	if err != nil {
		return
	}

	fmt.Println("received new update message")
	var entity data.Entity
	err = json.Unmarshal(event.Entity, &entity) //TODO: Error handling
	if err != nil {
		fmt.Println(err)
		return
	}
	entity.Name = s.ToLower(entity.Name)

	searchLocation, _ := client.GetGroups(ctx, token.AccessToken, realm, gocloak.GetGroupsParams{})
	if len(searchLocation) != 0 {
		hasUpdated := false
		for _, searchSubgroup := range *searchLocation[0].SubGroups {
			subgroup, err := client.GetGroup(ctx, token.AccessToken, realm, *searchSubgroup.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			locationIds := (*subgroup.Attributes)["locationId"]
			if len(locationIds) != 0 && locationIds[0] == strconv.Itoa(entity.Id) {
				*subgroup.Name = entity.Name
				err = client.UpdateGroup(ctx, token.AccessToken, realm, *subgroup)

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
