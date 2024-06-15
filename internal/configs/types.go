package configs

type (
	Config struct {
		Service  Service
		Database DatabaseConfig
		Redis    RedisConfig
	}

	Service struct {
		Port      string
		SecretKey string
	}

	DatabaseConfig struct {
		Master PostgresConfig
		Slave  PostgresConfig
	}

	PostgresConfig struct {
		Address string
	}

	RedisConfig struct {
		Address             string
		Password            string
		MaxActiveConnection int
		MaxIdleConnection   int
		TimeOut             int
		Wait                bool
		DB                  int
	}
)
