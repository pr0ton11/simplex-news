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
	client := client.NewSimplexClient("ws://localhost:3333", ctx)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.SendMessage("Marc", "Hello, World!")
	if err != nil {
		print(fmt.Sprintf("Error sending message: %v", err))
	}
	fmt.Println("Message sent successfully!")
}
