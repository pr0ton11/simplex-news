package main

import (
	"context"
	"fmt"
	"time"

	"github.com/proton11/simplex-news/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client := client.NewSimplexClient(ctx, "ws://localhost:3333")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
	err = client.ChangeDisplayName(fmt.Sprintf("Marc-%d", time.Now().Unix()))
	if err != nil {
		print(err)
	}
	defer client.Close()
	err = client.SendMessage("@Marc", "Hello, World!")
	if err != nil {
		print(err)
	}
}
