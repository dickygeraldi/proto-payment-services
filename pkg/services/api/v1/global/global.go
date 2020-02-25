package global

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Initializar .env variable
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
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

// Request to QREN
type Qren struct {
	MerchantApiKey string `json:"merchantApiKey"`
	Nominal        string `json:"nominal"`
	StaticQR       string `json:"staticQR"`
	InvoiceName    string `json:"invoiceName"`
	QrGaruda       string `json:"qrGaruda"`
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
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
