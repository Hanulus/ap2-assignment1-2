package main

import (
	"log"
	"net"
	"os"

	pb "github.com/Hanulus/ap2-generated/payment"
	"google.golang.org/grpc"
	"payment-service/internal/app"
	"payment-service/internal/repository"
	grpcTransport "payment-service/internal/transport/grpc"
	httpTransport "payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

func main() {
	// Connect to the database
	db, err := app.NewDB()
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	// Wire up layers: repository -> use case
	paymentRepo := repository.NewPostgresPaymentRepo(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)

	// Start gRPC server in a background goroutine
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9082"
	}
	go startGRPCServer(paymentUC, grpcPort)

	// Start REST server (kept for backward compatibility)
	handler := httpTransport.NewPaymentHandler(paymentUC)
	router := httpTransport.NewRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9081"
	}
	log.Printf("Payment REST server starting on :%s", port)
	log.Fatal(router.Run(":" + port))
}

func startGRPCServer(uc *usecase.PaymentUseCase, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port %s: %v", port, err)
	}

	// Register the logging interceptor (bonus: logs method + duration)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcTransport.LoggingInterceptor),
	)

	// Register our PaymentService implementation
	pb.RegisterPaymentServiceServer(grpcServer, grpcTransport.NewPaymentGRPCServer(uc))

	log.Printf("Payment gRPC server starting on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
