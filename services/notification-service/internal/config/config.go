package config

import "os"

type Config struct {
	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() Config {
	return Config{
		KafkaBrokers: getEnv(
			"KAFKA_BROKERS",
			"localhost:9092",
		),
		KafkaTopic: getEnv(
			"KAFKA_TOPIC",
			"order-events",
		),
		KafkaGroupID: getEnv(
			"KAFKA_GROUP_ID",
			"notification-service",
		),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
