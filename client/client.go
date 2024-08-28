package main

import (
	"context"
	"io"
	"log"
	"time"

	protobuff "github.com/ramasuryananda/grpc-learning/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	con, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to serve client  : %v", err)
	}

	defer con.Close()

	c := protobuff.NewCoffeShopClient(con)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &protobuff.MenuRequest{})
	if err != nil {
		log.Fatalf("failed to get menu  : %v", err)
	}

	done := make(chan bool)

	var items []*protobuff.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				log.Fatalf("failed failed receiving : %v", err)
			}

			items = resp.Items
			log.Printf("Resp received : %v", resp.Items)
		}
	}()

	<-done

	receipt, err := c.PlaceOrder(ctx, &protobuff.Order{Items: items})
	if err != nil {
		log.Fatalf("failed to get receipt  : %v", err)
	}
	log.Printf("%v", receipt)

	status, err := c.GetOrderStatus(ctx, &protobuff.Receipt{
		Id: receipt.Id,
	})
	if err != nil {
		log.Fatalf("failed to get status  : %v", err)
	}
	log.Printf("%v", status)

}
