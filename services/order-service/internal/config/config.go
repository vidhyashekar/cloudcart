package config

import "os"

type Config struct {
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	ProductServiceURL string
	KafkaBrokers      string
	KafkaTopic        string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8082"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "cloudcart"),
		JWTSecret:         getEnv("JWT_SECRET", "cloudcart-development-secret"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8081"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:        getEnv("KAFKA_TOPIC", "order-events"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
