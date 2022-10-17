package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

func RegisterDevice(deviceId string) string {
	ctx := context.Background()
	client := graphql.NewClient("http://localhost:6678/graphql", http.DefaultClient)
	resp, err := registerDevice(ctx, client, deviceId)

	fmt.Println(resp.RegisterDevice, err)
	return resp.RegisterDevice
}
