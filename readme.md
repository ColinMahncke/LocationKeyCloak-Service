## Location Keycloak Service

The service unmarshal kafka data and based on the data it creates groups and subgroups

## Environment variables

- host (config your own kafka host)
- keycloakhost (config your own keycloak host)
- user (config your keycloak user)
- password (config your keycloak user password)
- realmName (config your realmName)

## Function

The service loads the environment variables build a connection to the host and starts a kafka reader. After that the service reads the configurable variables for the keycloak login. The incoming data from Kafka are converted and the service creates new groups and subgroups in Keycloak (if they aren't already exist).