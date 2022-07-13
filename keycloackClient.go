package main

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v11"
)

func GetKeycloakClient() (client gocloak.GoCloak, ctx context.Context, token *gocloak.JWT, realm string, err error) {
	realm = Getenv("realmName", "realmName")
	client = gocloak.NewClient(Getenv("keycloakHost", "https://mycool.keycloak.instance"))
	ctx = context.Background()
	token, err = client.LoginAdmin(ctx, Getenv("user", "user"), Getenv("password", "password"), realm)
	if err != nil {
		fmt.Println(err)

	}
	return
}
