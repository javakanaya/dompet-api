package main

import (
	"dompet-api/config"
	"dompet-api/controller"
	"dompet-api/repository"
	"dompet-api/routes"
	"dompet-api/service"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetupDatabaseConnection()

	jwtService := service.NewJWTService()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, jwtService)

	server := gin.Default()

	routes.UserRouter(server, userController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9393"
	}
	server.Run(":" + port)

	config.CloseDatabaseConnection(db)
}
