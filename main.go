package main

import (
	"context"
	"log"
	"net"

	protobuff "github.com/ramasuryananda/grpc-learning/pb"
	"google.golang.org/grpc"
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

	for i := 0; i < len(items); i++ {
		srv.Send(&protobuff.Menu{
			Items: items[0 : i+1],
		})
	}

	return nil

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

	protobuff.RegisterCoffeShopServer(grpceServer, &server{})

	log.Printf("GRPC Server started")
	if err := grpceServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc  : %s", err)
	}

}
