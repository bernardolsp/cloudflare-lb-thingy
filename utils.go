package main

import (
	"fmt"
	"strings"
)

func CreatePayloadFromURLs(urls string) RequestBody {
	addresses := strings.Split(urls, ",")
	origins := make([]Origin, len(addresses))

	for i, address := range addresses {
		trimmedAddress := strings.TrimSpace(address)
		origins[i] = Origin{
			Address: trimmedAddress,
			Name:    fmt.Sprintf("origin-%d", i+1), // or use some part of the URL as name
		}
	}

	return RequestBody{Origins: origins, Name: "my-load-balancer"}
}
