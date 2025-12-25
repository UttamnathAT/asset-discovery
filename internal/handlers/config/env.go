package config

import (
	"os"
	"strconv"
	"time"
)

// env config
type AppEnv struct {
	Postgre PostgreEnv
	Server  ServerEnv
	Kafka   KafkaENv
	Auth    AuthEnv
}

type PostgreEnv struct {
	DNS PostgreDNSENv
	DB  PostgreDBEnv
}

type PostgreDNSENv struct {
	ReadWrite string
	ReadOnly  string
	Log       string
	Main      string
}

type PostgreDBEnv struct {
	Main string
	Log  string
}

type ServerEnv struct {
	GRPCPort     int
	HTTPPort     int
	ClientOrigin string
	Environment  string
	RedisAddr    string
	RedisDB      int
	IsLive       bool
	Smtp         SmtpEnv
}
type SmtpEnv struct {
	Host     string
	Port     int
	Email    string
	Password string
}
type KafkaENv struct {
	Addrs string
	Topic string
}

type AuthEnv struct {
	Secret string
}

func LoadEnv() (env AppEnv, err error) {
	env = AppEnv{
		Postgre: PostgreEnv{
			DNS: PostgreDNSENv{
				ReadWrite: getEnv("POSTGRE_WRITE_DNS", "host=localhost user=postgres password=postgres dbname=app_db port=5432 sslmode=disable"),
				ReadOnly:  getEnv("POSTGRE_READ_DNS", "host=localhost user=postgres password=postgres dbname=app_db port=5432 sslmode=disable"),
				Log:       getEnv("POSTGRE_LOG_DNS", "host=localhost user=postgres password=postgres dbname=log_db port=5432 sslmode=disable"),
				Main:      getEnv("POSTGRE_MAIN_DNS", "host=localhost user=postgres password=postgres port=5432 sslmode=disable"),
			},
			DB: PostgreDBEnv{
				Main: getEnv("POSTGRE_MAIN_DB", "app_db"),
				Log:  getEnv("POSTGRE_LOG_DB", "log_db"),
			},
		},
		Server: ServerEnv{
			GRPCPort:     getEnvInt("GRPC_POST", 50051),
			HTTPPort:     getEnvInt("HTTP_POST", 8003),
			ClientOrigin: getEnv("CLIENT_ORIGIN", "*"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
			RedisDB:      getEnvInt("REDIS_DB", 0),
			IsLive:       getEnvBool("IS_LIVE", false),
			Smtp: SmtpEnv{
				Host:     getEnv("SMTP_HOST", "smtp.example.com"),
				Port:     getEnvInt("SMTP_PORT", 587),
				Email:    getEnv("SMTP_EMAIL", "example@example.com"),
				Password: getEnv("SMTP_PASSWORD", "password"),
			},
		},
		Kafka: KafkaENv{
			Addrs: getEnv("KAFKA_ADDRS", "localhost:9092"),
			Topic: getEnv("KAFKA_TOPIC", "asset_discovery"),
		},
		Auth: AuthEnv{
			Secret: getEnv("JWT_SECRET", "secret"),
		},
	}
	return
}

// Helper functions to get environment variables with default values

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}
func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
