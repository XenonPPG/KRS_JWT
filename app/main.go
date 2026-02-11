package main

import (
	"JWT/internal/controllers"
	"JWT/internal/initializers"
	"JWT/internal/middleware"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {
	// grpc
	grpcConn, err := initializers.ConnectGRPC(os.Getenv("GRPC_ADDRESS"))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(grpcConn)

	// fiber
	fmt.Println("Starting server...")

	app := fiber.New()

	api := app.Group("/api")
	api.Route("/user", func(router fiber.Router) {
		router.Get("/:id", controllers.GetUser)
		router.Post("/", controllers.CreateUser)
		router.Get("/verify", controllers.VerifyPassword)

		protected := router.Group("/", middleware.JWTProtected)
		protected.Get("/", controllers.GetAllUsers)
		protected.Put("/", controllers.UpdateUser)
		protected.Put("/password", controllers.UpdatePassword)
		protected.Delete("/:id", controllers.DeleteUser)
	})

	err = app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		panic(err)
	}
}
