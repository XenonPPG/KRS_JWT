package main

import (
	"JWT/internal/controllers"
	"JWT/internal/initializers"
	"JWT/internal/middleware"
	"fmt"
	"os"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
	}))

	api := app.Group("/api")
	api.Route("/user", func(router fiber.Router) {
		// public
		router.Post("/", middleware.AutoLogout, controllers.CreateUser)
		router.Get("/:id", controllers.GetUser)

		// jwt protected
		protected := router.Group("", middleware.JWTProtected)
		protected.Put("/password", controllers.UpdatePassword)

		// targetID is defined by user role
		// if the role is not enough - id from params is ignored
		protected.Delete("/:id", controllers.DeleteUser)
		protected.Put("/:id", controllers.UpdateUser)

		// admins can control other accounts
		adminOnly := router.Group("", middleware.JWTProtected, middleware.RoleRequired(desc.UserRole_ADMIN))
		adminOnly.Get("/", controllers.GetAllUsers)
	})

	api.Route("/auth", func(router fiber.Router) {
		router.Post("/logout", controllers.LogoutHandler)
		router.Post("/refresh", controllers.RefreshTokens)
		router.Post("/login", middleware.AutoLogout, controllers.Login)
	})

	api.Route("/note", func(router fiber.Router) {
		router.Use(middleware.JWTProtected)

		router.Post("/", controllers.CreateNote)
		router.Get("/", controllers.GetAllNotes)
		router.Get("/:id", controllers.GetNote)
		router.Put("/", controllers.UpdateNote)
		router.Delete("/:id", controllers.DeleteNote)
	})

	api.Get("/health", controllers.HealthCheck)

	err = app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		panic(err)
	}
}
