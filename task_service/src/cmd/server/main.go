package main

import (
	"log"
	"net/http"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/interfaces/grpc/middleware"
	"task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/routes"
	"task_service/src/internal/usecase"
	logger "task_service/src/pkg/logger"

	"google.golang.org/grpc"
)

func main() {

	logger.Init()
	defer logger.Sync()

	database, err := persistance.ConnectToDatabase()
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database : %v", err)
	}

	sessionClient, err := sessionclient.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect session service: %v", err)
	}
	defer sessionClient.Close()

	_ = grpc.NewServer(
		grpc.UnaryInterceptor(middleware.SessionAuthInterceptor(sessionClient)),
	)

	// fmt.Println("grpc server : ", grpcServer)

	TaskRepo := persistance.NewTaskRepo(database)
	publisher := persistance.NewRedisPublisher("localhost:6379", "")
	TaskService := usecase.NewTaskService(TaskRepo)
	TaskHandler := handler.NewTaskHandler(TaskService, publisher)

	// router := routes.InitRoutes(*TaskHandler)
	router := routes.InitRoutes(*TaskHandler, sessionClient)

	// register your task service handlers here

	// ==== will define port instead of hard coded 80801 ====
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		logger.Log.Fatalf("Failed to start the server : %v", err)
	}

	logger.Log.Info("Server running on http://localhost:8080")
}
