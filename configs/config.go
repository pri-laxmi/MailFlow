package configs
import (
	"os"
	"log"
	"github.com/joho/godotenv"
)
type Config struct {
	Port string

	DBHost string
	DBPort string    
	DBUser string
	DBPassword string
	DBName string
	DBSSLMode string
}

func LoadConfig() *Config{
	err:=godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Port: os.Getenv("PORT"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSLMODE"),
	}
}