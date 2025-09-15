package main

import (
	"budget-book/internal/config"
	"budget-book/internal/infrastructure/database"
	infraRepo "budget-book/internal/infrastructure/repository"
	"budget-book/internal/interface/handler"
	"budget-book/internal/interface/middleware"
	"budget-book/internal/usecase"
	"log"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	cfg := config.Load()

	dbConfig := &database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	transactionRepo := infraRepo.NewTransactionRepository(db)
	categoryRepo := infraRepo.NewCategoryRepository(db)
	budgetRepo := infraRepo.NewBudgetRepository(db)

	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, categoryRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	budgetUseCase := usecase.NewBudgetUseCase(budgetRepo, categoryRepo)
	summaryUseCase := usecase.NewSummaryUseCase(transactionRepo, categoryRepo, budgetRepo)

	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	budgetHandler := handler.NewBudgetHandler(budgetUseCase)
	summaryHandler := handler.NewSummaryHandler(summaryUseCase)

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(echoMiddleware.Recover())

	api := e.Group("/api")

	api.GET("/transactions", transactionHandler.GetTransactions)
	api.POST("/transactions", transactionHandler.CreateTransaction)
	api.GET("/transactions/:id", transactionHandler.GetTransaction)
	api.PUT("/transactions/:id", transactionHandler.UpdateTransaction)
	api.DELETE("/transactions/:id", transactionHandler.DeleteTransaction)

	api.GET("/categories", categoryHandler.GetCategories)
	api.POST("/categories", categoryHandler.CreateCategory)
	api.GET("/categories/:id", categoryHandler.GetCategory)
	api.PUT("/categories/:id", categoryHandler.UpdateCategory)
	api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	api.GET("/budgets", budgetHandler.GetBudgets)
	api.POST("/budgets", budgetHandler.CreateBudget)
	api.GET("/budgets/:id", budgetHandler.GetBudget)
	api.PUT("/budgets/:id", budgetHandler.UpdateBudget)
	api.DELETE("/budgets/:id", budgetHandler.DeleteBudget)

	api.GET("/summary/:year/:month", summaryHandler.GetMonthlySummary)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(e.Start(":" + cfg.Server.Port))
}
