package main

import (
	"dompet-api/config"
	"dompet-api/controller"
	"dompet-api/middleware"
	"dompet-api/repository"
	"dompet-api/routes"
	"dompet-api/service"
	"log"
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

	defer config.CloseDatabaseConnection(db)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRouter(server, userController, jwtService)
	routes.DompetRouter(server, dompetController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
