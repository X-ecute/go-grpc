package grpc

import (
	"context"
	"log"
	"net"

	"github.com/X-ecute/go-grpc/internal/rocket"
	rkt "github.com/X-ecute/go-grpc/protos/rocket/v1/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // Add this for debugging
)

type RocketService interface {
	GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error)
	InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error)
	DeleteRocket(ctx context.Context, id string) error
}

// Handler - will handle new gRPC requests
type Handler struct {
	rkt.UnimplementedRocketServiceServer
	RocketService RocketService
}

func New(rktService RocketService) *Handler { // Return pointer
	return &Handler{
		RocketService: rktService,
	}
}

func (h *Handler) Serve() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	rkt.RegisterRocketServiceServer(grpcServer, h)

	// Add reflection for easier debugging
	reflection.Register(grpcServer)

	log.Println("gRPC server starting on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}

// ✅ FIXED: Method names must match .proto definition
func (h *Handler) GetRocket(ctx context.Context, req *rkt.GetRocketRequest) (*rkt.GetRocketResponse, error) {
	log.Printf("GetRocket called with: %v", req)

	// Use your business logic service
	rocketByID, err := h.RocketService.GetRocketByID(ctx, req.Id)
	if err != nil {
		log.Printf("GetRocket failed with: %v", err)
		return nil, err
	}

	return &rkt.GetRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   rocketByID.ID,
			Name: rocketByID.Name,
			Type: rocketByID.Type,
		},
	}, nil
}

// ✅ Keep this name as it matches .proto
func (h *Handler) AddRocket(ctx context.Context, req *rkt.AddRocketRequest) (*rkt.AddRocketResponse, error) {
	log.Printf("AddRocket called with: %v", req)

	rocketToInsert := rocket.Rocket{
		ID:   req.Rocket.Id,
		Name: req.Rocket.Name,
		Type: req.Rocket.Type,
	}

	insertedRocket, err := h.RocketService.InsertRocket(ctx, rocketToInsert)
	if err != nil {
		return nil, err
	}

	return &rkt.AddRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   insertedRocket.ID,
			Name: insertedRocket.Name,
			Type: insertedRocket.Type,
		},
	}, nil
}

// ✅ Keep this name as it matches .proto
func (h *Handler) DeleteRocket(ctx context.Context, req *rkt.DeleteRocketRequest) (*rkt.DeleteRocketResponse, error) {
	log.Printf("DeleteRocket called with: %v", req)

	err := h.RocketService.DeleteRocket(ctx, req.Rocket.Id)
	if err != nil {
		return &rkt.DeleteRocketResponse{
			Status: "false",
		}, err
	}

	return &rkt.DeleteRocketResponse{
		Status: "true",
	}, nil
}
