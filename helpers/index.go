package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

func EnvVar(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}
