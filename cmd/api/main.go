package main

import (
	"app/myproj/internal/api"
	"app/myproj/internal/auth"
	"app/myproj/internal/order"
	"app/myproj/pkg/workflow"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database setup
	dsn := "host=localhost user=postgres password=postgres dbname=auth port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate database
	db.AutoMigrate(&auth.User{}, &order.Order{}, &workflow.WorkflowInstance{})

	// Initialize services
	authService := auth.NewAuthService(db)
	orderService := order.NewOrderService(db)

	// Initialize workflow engine
	engine := workflow.NewWorkflowEngine(db)

	// Register workflows
	engine.RegisterWorkflow(auth.NewAuthWorkflow(authService))
	engine.RegisterWorkflow(order.NewOrderWorkflow(orderService))

	// Initialize handlers
	authHandler := api.NewAuthHandler(engine, authService)
	orderHandler := api.NewOrderHandler(engine, orderService)

	// Setup router
	router := gin.Default()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/start", authHandler.StartAuth)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/set-pin", authHandler.SetPIN)
		auth.POST("/verify-pin", authHandler.VerifyPIN)
	}

	// Order routes
	orders := router.Group("/orders")
	{
		orders.POST("/", orderHandler.CreateOrder)
		orders.POST("/:id/payment", orderHandler.ConfirmPayment)
		orders.POST("/:id/fulfill", orderHandler.FulfillOrder)
	}

	router.Run(":8080")
}
