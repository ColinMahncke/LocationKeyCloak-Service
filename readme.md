# Location Keycloak Service

The service unmarshal kafka data and based on the data it creates/deletes groups and subgroups
___

## Environment variables

- host (config your own kafka host)
- keycloakhost (config your own keycloak host)
- user (config your keycloak user)
- password (config your keycloak user password)
- realmName (config your realmName)
___

## Function

The service loads the environment variables build a connection to the host and starts a kafka reader. After that the service reads the configurable variables for the keycloak login. The incoming data from Kafka are converted and the service creates new groups and subgroups in Keycloak (if they aren't already exist). The Service can also delete groups based on the incoming kafka data (if there are matching groups).
___
## Kafka events


### Delete event
```json
{
	"EventType": "Delete",
	"Entity": {
		"id": 3,
		"name": "hamburg",
		"shortname": "commodo consequat",
		"geometry": {
			"type": "Point",
			"coordinates": [
				-31512158.275429994,
				-25747854.676155195
			]
		}
	}
}
```

### Create event
```json
{
	"EventType": "Create",
	"Entity": {
		"id": 3,
		"name": "hamburg",
		"shortname": "commodo consequat",
		"geometry": {
			"type": "Point",
			"coordinates": [
				-31512158.275429994,
				-25747854.676155195
			]
		}
	}
}
```
