package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	protobuff "github.com/ramasuryananda/grpc-learning/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFloat generates a random float32 between min and max
func randomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

type server struct {
	protobuff.UnimplementedCoffeShopServer
}

func (s *server) GetMenu(menuRequest *protobuff.MenuRequest, srv protobuff.CoffeShop_GetMenuServer) error {
	log.Printf("Received request to get menu")
	items := []*protobuff.Item{
		{
			Id:   "1",
			Name: "Coffee",
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
			log.Printf("Streaming canceled")
			return status.Error(codes.Canceled, "Stream ended")
		default:
			time.Sleep(2 * time.Second)
			index := i % len(items)
			i = i + 1

			log.Printf("Sending menu item: %v", items[index])

			err := srv.Send(&protobuff.Menu{
				Items: []*protobuff.Item{
					items[index],
				},
			})
			if err != nil {
				return status.Error(codes.Aborted, "Failed sending message")
			}
		}
	}
}

func (s *server) GetSahamData(menuRequest *protobuff.MenuRequest, srv protobuff.CoffeShop_GetSahamDataServer) error {
	log.Printf("Received request to get saham data")

	for {
		select {
		case <-srv.Context().Done():
			log.Printf("Streaming canceled")
			return status.Error(codes.Canceled, "Stream ended")
		default:
			time.Sleep(2 * time.Second)

			sahamData := &protobuff.Saham{
				Date:   time.Now().Format("2006-01-02 15:04:05"),
				Open:   randomFloat(100, 200),
				High:   randomFloat(100, 200),
				Low:    randomFloat(100, 200),
				Close:  randomFloat(100, 200),
				Volume: 12345678,
			}

			log.Printf("Sending saham data: %v", sahamData)

			if err := srv.Send(sahamData); err != nil {
				return status.Error(codes.Aborted, "Failed sending message")
			}
		}
	}
}

func (s *server) PlaceOrder(ctx context.Context, od *protobuff.Order) (*protobuff.Receipt, error) {
	log.Printf("Received request to place order")
	return &protobuff.Receipt{
		Id: "123",
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *protobuff.Receipt) (*protobuff.OrderStatus, error) {
	log.Printf("Received request to get order status")
	return &protobuff.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001: %v", err)
	}

	grpcServer := grpc.NewServer()
	serverService := &server{}

	protobuff.RegisterCoffeShopServer(grpcServer, serverService)

	log.Printf("gRPC Server started on port 9001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
