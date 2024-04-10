package astparser

import (
	"fmt"
	"regexp"
	"strings"
)

type ValueFinder interface {
	GetValue(string) string
}

// ReValueFinder возвращает ValueFinder, который ищет значение
// в строке по регулярному выражению `tagname:"(.+?)"` и возвращает
// первую группу, то есть все содержимое тега между двойными кавычками
func ReValueFinder(tagname string) ValueFinder {
	return &reValueFinder{re: regexp.MustCompile(fmt.Sprintf(`%s:"(.+?)"`, tagname))}
}

// MatchStringFinder возвращает ValueFinder, который возвращает всю
// строку, если strings.Contains(s, tagname)
func MatchStringFinder(tagname string) ValueFinder {
	return &matchStringFinder{tagname: tagname}
}

type reValueFinder struct {
	re *regexp.Regexp
}

func (v *reValueFinder) GetValue(s string) (value string) {
	var out string
	if v.re.MatchString(s) {
		out = v.re.FindStringSubmatch(s)[1]
	}
	return out
}

type matchStringFinder struct {
	tagname string
}

func (m *matchStringFinder) GetValue(s string) (value string) {
	searchStr := fmt.Sprintf("%s:", m.tagname)
	if !strings.Contains(s, searchStr) {
		return ""
	}
	return s
}
