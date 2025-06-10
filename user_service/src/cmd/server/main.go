package main

import (
	"log"
	"net"
	"net/http"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/interfaces/grpc/server"
	"user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/routes"
	"user_service/src/internal/usecase"
	logger "user_service/src/pkg/logger"

	pb "user_service/src/internal/interfaces/grpc/generated"

	"google.golang.org/grpc"
)

func main() {

	logger.Init()
	defer logger.Sync()

	database, err := persistance.ConnectToDatabase()
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database : %v", err)
	}

	UserRepo := persistance.NewUserRepo(database)
	UserService := usecase.NewUserService(UserRepo)
	UserHandler := handler.NewUserHandler(UserService)

	router := routes.InitRoutes(*UserHandler)

	// redisClient := redis.NewClient()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	sessionServer := server.NewSessionServer(persistance.RedisClient)
	pb.RegisterSessionValidatorServer(grpcServer, sessionServer)

	go func() {
		log.Println("gRPC server listening on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// ==== will define port instead of hard coded 8080 ====
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Log.Fatalf("Failed to start the server : %v", err)
	}

	logger.Log.Info("Server running on http://localhost:8080")
}
