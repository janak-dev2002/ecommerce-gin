package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	AccessTokenMinutes int
	RefreshTokenDays   int
	S3Endpoint         string `mapstructure:"S3_ENDPOINT"`
	S3Key              string `mapstructure:"S3_KEY"`
	S3Secret           string `mapstructure:"S3_SECRET"`
	S3Bucket           string `mapstructure:"S3_BUCKET"`
	S3Region           string `mapstructure:"S3_REGION"`
	S3PublicURL        string `mapstructure:"S3_PUBLIC_URL"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisPassword      string `mapstructure:"REDIS_PASSWORD"`
}

var Cfg Config

func LoadConfig() {
	// load .env if present (works locally; in prod rely on env)
	_ = godotenv.Load()

	accessMin, err := strconv.Atoi(getEnv("ACCESS_TOKEN_MINUTES", "15"))
	if err != nil {
		log.Println("Invalid ACCESS_TOKEN_MINUTES, using 15")
		accessMin = 15
	}
	refreshDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_DAYS", "7"))
	if err != nil {
		log.Println("Invalid REFRESH_TOKEN_DAYS, using 7")
		refreshDays = 7
	}

	Cfg = Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "ecommerce"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret"),
		AccessTokenMinutes: accessMin,
		RefreshTokenDays:   refreshDays,
		S3Endpoint:         getEnv("S3_ENDPOINT", ""),
		S3Key:              getEnv("S3_KEY", ""),
		S3Secret:           getEnv("S3_SECRET", ""),
		S3Bucket:           getEnv("S3_BUCKET", ""),
		S3Region:           getEnv("S3_REGION", "auto"),
		S3PublicURL:        getEnv("S3_PUBLIC_URL", ""),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
	}
	log.Println("Config loaded")
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultValue
}
