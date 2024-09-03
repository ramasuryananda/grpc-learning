package main

import (
	"context"
	"io"
	"log"

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

	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	menuStream, err := c.GetMenu(ctx, &protobuff.MenuRequest{})
	if err != nil {
		log.Fatalf("failed to get menu  : %v", err)
	}

	done := make(chan bool)

	var items []*protobuff.Item

	i := 0

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
			if i > 10 {
				cancel()
			}

			i++
		}
	}()

	<-done

	receipt, err := c.PlaceOrder(context.Background(), &protobuff.Order{Items: items})
	if err != nil {
		log.Fatalf("failed to get receipt  : %v", err)
	}
	log.Printf("%v", receipt)

	status, err := c.GetOrderStatus(context.Background(), &protobuff.Receipt{
		Id: receipt.Id,
	})
	if err != nil {
		log.Fatalf("failed to get status  : %v", err)
	}
	log.Printf("%v", status)

}
