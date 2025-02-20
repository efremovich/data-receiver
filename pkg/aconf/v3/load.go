package aconf

import (
	"context"
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

var (
	ErrNotAPointer = errors.New("переданный аргумент не является указателем на структуру")
	ErrFailedLoad  = errors.New("не удалось загрузить переменные окружения")
	ErrValidate    = errors.New("ошибка валидации структуры конфигурации")
)

// Load принимает в качестве аргумента указатель на структуру конфигурации
// и ищет переменные окружения, указанные в тегах "env:" полей структуры.
// Затем производится валидация полей структуры с помощью библиотеки go-playground/validator
//
// ВАЖНО: структуру нужно передавать по указателю, иначе метод
// вернет ошибку. Детали работы метода в README.md.
func Load(v any) error {
	return doLoadAndValidate(v)
}

// PreloadEnvsFile - загружает указанный аргументом path файл,
// считывает из него переменные окружения и задает их для текущей
// среды исполнения.
//
// ВАЖНО: Если переменная окружения уже задана, она не будет перезаписана вызовом
// этой функции.
func PreloadEnvsFile(path string) error {
	return preloadEnvs(path)
}

func doLoadAndValidate(v interface{}) error {
	if !isPtrToStruct(v) {
		return ErrNotAPointer
	}

	if err := envconfig.Process(context.Background(), v); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedLoad, err.Error())
	}
	return validateStruct(v)
}

func preloadEnvs(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedLoad, err.Error())
	}
	return nil
}
