package configs

type (
	Config struct {
		Service  Service
		Database DatabaseConfig
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
)
