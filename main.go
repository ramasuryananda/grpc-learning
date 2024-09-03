package main

import (
	"context"
	"log"
	"net"
	"time"

	protobuff "github.com/ramasuryananda/grpc-learning/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	protobuff.UnimplementedCoffeShopServer
}

func (s *server) GetMenu(menuRequest *protobuff.MenuRequest, srv protobuff.CoffeShop_GetMenuServer) error {
	log.Printf("received message to get menu")
	items := []*protobuff.Item{
		{
			Id:   "1",
			Name: "Coffe",
		},
		{
			Id:   "2",
			Name: "Tea",
		},
		{
			Id:   "3",
			Name: "Milk Tea",
		},
	}
	i := 0

	for {
		select {
		case <-srv.Context().Done():
			log.Printf("streaming canceled")
			return status.Error(codes.Canceled, "stream ended")
		default:
			time.Sleep(2 * time.Second)
			index := i % 3
			i = i + 1

			log.Printf("sending message get menu : %v", items[index])

			err := srv.Send(&protobuff.Menu{
				Items: []*protobuff.Item{
					items[index],
				},
			})
			if err != nil {
				return status.Error(codes.Aborted, "failed sending message")
			}
		}
	}
}

func (s *server) PlaceOrder(ctx context.Context, od *protobuff.Order) (*protobuff.Receipt, error) {
	log.Printf("received message to place order")
	return &protobuff.Receipt{
		Id: "123",
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *protobuff.Receipt) (*protobuff.OrderStatus, error) {
	log.Printf("received message to get order status")
	return &protobuff.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen to port : %s", err)
	}

	grpceServer := grpc.NewServer()
	serverService := &server{}

	protobuff.RegisterCoffeShopServer(grpceServer, serverService)

	log.Printf("GRPC Server started")
	if err := grpceServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc  : %s", err)
	}

}
