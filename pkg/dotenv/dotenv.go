package dotenv

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	kwEnvName   = "KW_ENV"
	Production  = "prod"
	Development = "dev"
	Test        = "test"
	Local       = "local"
)

func Load() {
	env := os.Getenv(kwEnvName)
	if env == "" {
		env = Development
	}

	// Пробуем загрузить локальную версию окружения
	if Test != env {
		_ = godotenv.Load(".env." + Local)
	}

	// Загружаем в зависимости от окружения
	_ = godotenv.Load(".env." + env)

	// Загружаем дефолтную .env
	_ = godotenv.Load()
}

func GetCurrentEnv() string {
	if envVal, ok := os.LookupEnv(kwEnvName); ok {
		return envVal
	}

	return Development
}

func GetString(varName string, defaultValue string) string {
	if envVal, ok := os.LookupEnv(varName); ok {
		return envVal
	}

	return defaultValue
}

func GetInt(varName string, defaultValue int) int {
	envVal, ok := os.LookupEnv(varName)

	if !ok {
		return defaultValue
	}

	if num, err := strconv.Atoi(envVal); err == nil {
		return num
	}

	return defaultValue
}
