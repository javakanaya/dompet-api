package main

import (
	"log"
	"dompet-api/config"
	"dompet-api/controller"
	"dompet-api/middleware"
	"dompet-api/repository"
	"dompet-api/routes"
	"dompet-api/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	db := config.SetupDatabaseConnection()

	jwtService := service.NewJWTService()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, jwtService)

	dompetRepository := repository.NewDompetRepository(db)
	dompetService := service.NewDompetService(dompetRepository)
	dompetController := controller.NewDompetController(dompetService)

	// mereka saling dependen
	defer config.CloseDatabaseConnection(db)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRouter(server, userController, dompetController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}

