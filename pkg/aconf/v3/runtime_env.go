package aconf

import (
	"os"
	"strings"
)

// GoEnvironment - переменная окружения, которая
// содержит среду запуска приложения.
const GoEnvironment string = "GO_ENVIRONMENT"

const (
	LOCAL       string = "local"
	DEVELOPMENT string = "development"
	STAGING     string = "staging"
	PRODUCTION  string = "production"
)

// IsProduction - проверяет, что приложение запущено в среде Production
func IsProduction() bool {
	return envEquals(PRODUCTION)
}

// IsStaging - проверяет, что приложение запущено в среде Staging
func IsStaging() bool {
	return envEquals(STAGING)
}

// IsDevelopment - проверяет, что приложение запущено в среде Development
func IsDevelopment() bool {
	return envEquals(DEVELOPMENT)
}

// IsLocal - проверяет, что приложение запущено в среде "local"
func IsLocal() bool {
	return envEquals(LOCAL)
}

// GetEnvironment - возвращает значение среды выполнения. Производится
// поиск переменной окружения, заданной константой GoEnvironment.
// Если такая переменная окружения не задана или ее значение не соответствует
// одному из известных, то вернется LOCAL.
func GetEnvironment() (env string) {
	return getEnv()
}

// getEnv вернет найденное значение среды. Если
// переменная заданная константой goEnvironment пуста,
// вернет константу DEVELOPMENT
func getEnv() (env string) {
	found := strings.ToLower(os.Getenv(GoEnvironment))
	if !envValid(found) {
		return LOCAL
	}
	return found
}

func envValid(env string) bool {
	var valid bool
	for _, v := range []string{LOCAL, DEVELOPMENT, STAGING, PRODUCTION} {
		if v == env {
			valid = true
			break
		}
	}
	return valid
}

func envEquals(target string) bool {
	env := getEnv()
	return env == target
}
