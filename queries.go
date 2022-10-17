package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

func GetCollectionLatest(collectionId string) string {
	ctx := context.Background()
	client := graphql.NewClient("http://localhost:6678/graphql", http.DefaultClient)
	resp, err := collectionLatest(ctx, client, collectionId)

	fmt.Println(resp.GetCollectionLatest().UnsplashUrl, err)
	return resp.GetCollectionLatest().UnsplashUrl
}
