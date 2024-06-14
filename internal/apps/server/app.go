package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/internal/handler/books"
	"github.com/yeremiaaryo/gotu-assignment/internal/handler/orders"
	"github.com/yeremiaaryo/gotu-assignment/internal/handler/users"
	auth "github.com/yeremiaaryo/gotu-assignment/internal/middleware"
	booksRepository "github.com/yeremiaaryo/gotu-assignment/internal/repository/books"
	ordersRepository "github.com/yeremiaaryo/gotu-assignment/internal/repository/orders"
	usersRepository "github.com/yeremiaaryo/gotu-assignment/internal/repository/users"
	booksUsecase "github.com/yeremiaaryo/gotu-assignment/internal/usecase/books"
	ordersUsecase "github.com/yeremiaaryo/gotu-assignment/internal/usecase/orders"
	usersUsecase "github.com/yeremiaaryo/gotu-assignment/internal/usecase/users"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitApps(cfg *configs.Config) error {
	// Database initialization
	masterDB, err := internalsql.OpenMasterDB("postgres", cfg.Database.Master.Address)
	if err != nil {
		log.Fatalf("cannot open master db connection, err: %v", err)
		return err
	}
	slaveDB, err := internalsql.OpenSlaveDB("postgres", cfg.Database.Slave.Address)
	if err != nil {
		log.Fatalf("cannot open slave db connection, err: %v", err)
		return err
	}

	// Init all repo here
	usersRepo := usersRepository.New(masterDB, slaveDB)
	booksRepo := booksRepository.New(masterDB, slaveDB)
	ordersRepo := ordersRepository.New(masterDB, slaveDB)

	// Init all usecase here
	usersUsecase := usersUsecase.New(usersRepo, cfg)
	booksUsecase := booksUsecase.New(booksRepo, cfg)
	ordersUsecase := ordersUsecase.New(ordersRepo, booksRepo, cfg)

	// Init all handler here
	usersHandler := users.New(usersUsecase)
	booksHandler := books.New(booksUsecase)
	ordersHandler := orders.New(ordersUsecase)

	// Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// User handler
	e.POST("/register", usersHandler.CreateUser)
	e.POST("/login", usersHandler.Login)

	// Book handler
	e.GET("/books", booksHandler.GetBooks)

	// Order handler
	e.POST("/order", ordersHandler.CreateOrder, auth.AuthMiddleware)

	// Start server
	e.Logger.Fatal(e.Start(cfg.Service.Port))
	return nil
}
