package global

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// Initializar .env variable
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

// Tokenization for JWT auth
type MerchantData struct {
	MerchantId string
	Password   string
	jwt.StandardClaims
}

// Struct message error
type MessageError struct {
	Code    string
	Message string
}

// structure redis configuration
type PostGresql struct {
	PostgreClient string
	PostgrePass   string
	PosgreUser    string
	PostgrePort   string
	PostgreDb     string
}

// structure all configuration
type Configuration struct {
	Postgresql PostGresql
	ApiVersion string
	Token      string
	KeyPass    string
	KeyAes     string
}

// New returns a new Config struct
func New() *Configuration {
	return &Configuration{
		Postgresql: PostGresql{
			PostgreClient: getEnv("POSTGRESQL_URL", ""),
			PostgrePass:   getEnv("POSTGRESQL_PASS", ""),
			PostgreDb:     getEnv("POSTGRESQL_DB", ""),
			PosgreUser:    getEnv("POSTGRESQL_USER", ""),
			PostgrePort:   getEnv("POSTGRESQL_PORT", ""),
		},
		ApiVersion: getEnv("API_VERSION", ""),
		Token:      getEnv("TOKEN", ""),
		KeyPass:    getEnv("KEY_PASS", ""),
		KeyAes:     getEnv("KEY_AES", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// func to set code and message validation
func GetMessageError() map[int]MessageError {
	dataError := make(map[int]MessageError)

	dataError[00] = MessageError{"00", "Processing Success"}
	dataError[422] = MessageError{"422", "Processing Data Error"}
	dataError[01] = MessageError{"001", "Processing pending, try again later"}

	return dataError
}
