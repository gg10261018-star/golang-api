package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load()
}

func Getenv(key string, fallback ...string) string {
	v := os.Getenv(key)
	if v == "" && len(fallback) > 0 {
		return fallback[0]
	}
	return v
}

func GetenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
