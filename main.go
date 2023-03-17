package main

import (
	"log"
	"oprec/dompet-api/config"
	"oprec/dompet-api/controller"
	"oprec/dompet-api/middleware"
	"oprec/dompet-api/repository"
	"oprec/dompet-api/routes"
	"oprec/dompet-api/service"
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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	dompetRepository := repository.NewDompetRepository(db)
	dompetService := service.NewDompetService(dompetRepository)
	dompetController := controller.NewDompetController(dompetService)

	// mereka saling dependen
	defer config.CloseDatabaseConnection(db)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRoutes(server, userController, dompetController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
