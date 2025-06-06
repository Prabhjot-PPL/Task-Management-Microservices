package main

import (
	"net/http"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/interfaces/input/api/rest/handler"
	"user_service/src/internal/interfaces/input/api/rest/routes"
	"user_service/src/internal/usecase"
	logger "user_service/src/pkg/logger"
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

	// ==== will define port instead of hard coded 8080 ====
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Log.Fatalf("Failed to start the server : %v", err)
	}

	logger.Log.Info("Server running on http://localhost:8080")
}
