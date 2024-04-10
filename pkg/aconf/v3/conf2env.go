package aconf

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/efremovich/data-receiver/pkg/aconf/v3/internal/astparser"
)

const tagnameEnv string = "env"

// grpEndlOrKEyword - группа регулярки, означающая конец строки или следующее ключевое слово
const grpEndlOrKeyword = `($|decodeunset|noinit|separator|delimiter|overwrite|prefix|required|default)`

var (
	reEnvName    = regexp.MustCompile(` *([\w-_]+)`) // ищет первое слово в строке, состоящее из [a-zA-Z0-9_-]
	rePrefixName = regexp.MustCompile(`prefix= *([\w-_]+)`)
	reDefault    = regexp.MustCompile(fmt.Sprintf(`default= *(.+?)[, ]*%s`, grpEndlOrKeyword))
)

// ConfFileToEnvSlice принимает имя файла исходного текста на Go и имя структуры
// по тегам которой нужно сгенерировать переменные окружения для конфигурации. Возвращает
// массив строк, в которых переменным окружения присвоены значения из тега env:, указанные в default,
// например для `env:"MY_ENV, default=1"` вернется строка "MY_ENV=1".
// Если default не указан, то вернется строка, указывающая на тип, который должен быть
// задан для переменной, например "MY_ENV=int".
func ConfFileToEnvSlice(filename, structname string) ([]string, error) {
	parser, err := astparser.NewAstParser(filename)
	if err != nil {
		return nil, err
	}
	finder := astparser.ReValueFinder(tagnameEnv)
	out, err := parser.FindTagValuesInStructDecl(structname, finder)
	if err != nil {
		return nil, err
	}

	result, err := mapToEnvsStringSlice(out)
	if err != nil {
		return nil, err
	}

	sort.Strings(result)
	return result, nil
}

// mapToEnvsStringSlice конвертирует map[string]any, где срока - содержимое тега
// env, а значение - строковое представление базового типа, либо вложенная map[string]any
// если полем структуры была другая структура.
func mapToEnvsStringSlice(m map[string]any) ([]string, error) {
	var out []string
	for k, v := range m {
		if inner_m, ok := v.(map[string]any); ok {
			inner, err := mapToEnvsStringSlice(inner_m)
			if err != nil {
				return nil, err
			}
			prefix := findPrefix(k)
			for _, s := range inner {
				s = fmt.Sprintf("%s%s", prefix, s)
				out = append(out, s)
			}
			continue
		}

		if val, ok := v.(string); ok {
			name, err := findName(k)
			if err != nil {
				return nil, err
			}
			s := fmt.Sprintf("%s=%v", name, defaultOrExample(k, val))
			out = append(out, s)
			continue
		}

		return nil, fmt.Errorf("в переданной структуре значение не является map[string]any или string: %+v", m)
	}
	return out, nil
}

func findName(s string) (string, error) {
	if !reEnvName.MatchString(s) {
		return "", fmt.Errorf("для тэга \"%s:\" не найдено имя переменной окружения, нужно указать в тэге `%s:\"SOME_NAME\"`, исходная строка: %s", tagnameEnv, tagnameEnv, s)
	}
	return reEnvName.FindStringSubmatch(s)[1], nil
}

func findPrefix(s string) string {
	if !rePrefixName.MatchString(s) {
		return ""

	}
	return rePrefixName.FindStringSubmatch(s)[1]
}

func defaultOrExample(tag, val string) string {
	if !reDefault.MatchString(tag) {
		return val
	}
	def := reDefault.FindStringSubmatch(tag)[1]
	if strings.Contains(def, ",") || strings.Contains(def, ":") {
		// возьмем в одинарные кавычки значение генерируемой переменой
		def = fmt.Sprintf("'%s'", def)
	}
	return def
}
