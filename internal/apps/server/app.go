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
	"github.com/yeremiaaryo/gotu-assignment/pkg/redis"
	"log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitApps(cfg *configs.Config) error {
	redisAgent, err := initRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("init redis failed: %v", err)
	}

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
	booksRepo := booksRepository.New(masterDB, slaveDB, redisAgent)
	ordersRepo := ordersRepository.New(masterDB, slaveDB)

	// Init all usecase here
	usersUsecase := usersUsecase.New(usersRepo, redisAgent, cfg)
	booksUsecase := booksUsecase.New(booksRepo, cfg)
	ordersUsecase := ordersUsecase.New(ordersRepo, booksRepo, cfg)

	// Init all handler here
	usersHandler := users.New(usersUsecase)
	booksHandler := books.New(booksUsecase)
	ordersHandler := orders.New(ordersUsecase)

	// init auth
	authHandler := auth.New(redisAgent)

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
	e.POST("/order", ordersHandler.CreateOrder, authHandler.AuthMiddleware)
	e.GET("/order", ordersHandler.GetOrderHistory, authHandler.AuthMiddleware)

	// Start server
	e.Logger.Fatal(e.Start(cfg.Service.Port))
	return nil
}

func initRedis(config *configs.RedisConfig) (*redis.Redis, error) {
	// init redis MS configs.
	rdsConfig := redis.RedisConfig{
		Address:  config.Address,
		Password: config.Password,
		Options: []redis.RedisOptions{
			{
				MaxActive: config.MaxActiveConnection,
				MaxIdle:   config.MaxIdleConnection,
				Timeout:   config.TimeOut,
				Wait:      config.Wait,
			},
		},
	}

	// get MS redis agent.
	redisAgent := redis.NewRedis(rdsConfig)
	err := redisAgent.Ping()
	if err != nil {
		return nil, err
	}

	return redisAgent, nil
}
