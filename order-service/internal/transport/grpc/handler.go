package grpc

import (
	"time"

	pb "github.com/Hanulus/ap2-generated/orderstream"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order-service/internal/usecase"
)

// OrderGRPCServer implements the gRPC OrderService (streaming)
type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderGRPCServer(uc *usecase.OrderUseCase) *OrderGRPCServer {
	return &OrderGRPCServer{uc: uc}
}

// SubscribeToOrderUpdates streams order status changes to the client.
// It polls the DB every second and sends an update when the status changes.
func (s *OrderGRPCServer) SubscribeToOrderUpdates(
	req *pb.OrderRequest,
	stream pb.OrderService_SubscribeToOrderUpdatesServer,
) error {
	if req.OrderId == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	lastStatus := ""

	// Poll DB until client disconnects or context is cancelled
	for {
		// Check if client disconnected
		if stream.Context().Err() != nil {
			return nil
		}

		// Fetch current order from DB
		order, err := s.uc.GetOrder(req.OrderId)
		if err != nil {
			return status.Error(codes.NotFound, "order not found")
		}

		// Only send update if status actually changed
		if order.Status != lastStatus {
			lastStatus = order.Status

			err = stream.Send(&pb.OrderStatusUpdate{
				OrderId: order.ID,
				Status:  order.Status,
			})
			if err != nil {
				return err
			}
		}

		// Wait 1 second before next DB check
		time.Sleep(1 * time.Second)
	}
}
