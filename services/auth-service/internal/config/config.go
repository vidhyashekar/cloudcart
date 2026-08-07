package config

import "os"

// Config holds the configuration values for the application.
type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load loads the configuration from environment variables and returns a Config struct.
func Load() *Config {

	return &Config{
		Port: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "cloudcart"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is empty or not present, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
