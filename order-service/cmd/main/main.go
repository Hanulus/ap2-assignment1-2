package main

import (
	"log"
	"net"
	"os"

	pbStream "github.com/Hanulus/ap2-generated/orderstream"
	"google.golang.org/grpc"
	"order-service/internal/app"
	"order-service/internal/repository"
	grpcTransport "order-service/internal/transport/grpc"
	httpTransport "order-service/internal/transport/http"
	"order-service/internal/usecase"
)

func main() {
	// Connect to the database
	db, err := app.NewDB()
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	// gRPC address of the Payment Service (from env, not hardcoded)
	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentGRPCAddr == "" {
		paymentGRPCAddr = "localhost:9082"
	}

	// Use gRPC client instead of the old REST client
	paymentClient, err := repository.NewPaymentGRPCClient(paymentGRPCAddr)
	if err != nil {
		log.Fatalf("failed to connect to payment gRPC: %v", err)
	}

	// Wire up layers
	orderRepo := repository.NewPostgresOrderRepo(db)
	orderUC := usecase.NewOrderUseCase(orderRepo, paymentClient)

	// Start order streaming gRPC server in background
	streamPort := os.Getenv("GRPC_PORT")
	if streamPort == "" {
		streamPort = "9083"
	}
	go startStreamingServer(orderUC, streamPort)

	// Start REST server (external API stays REST)
	handler := httpTransport.NewOrderHandler(orderUC)
	router := httpTransport.NewRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}
	log.Printf("Order REST server starting on :%s", port)
	log.Fatal(router.Run(":" + port))
}

func startStreamingServer(uc *usecase.OrderUseCase, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on streaming port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Register the streaming OrderService
	pbStream.RegisterOrderServiceServer(grpcServer, grpcTransport.NewOrderGRPCServer(uc))

	log.Printf("Order streaming gRPC server starting on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("streaming gRPC server error: %v", err)
	}
}
