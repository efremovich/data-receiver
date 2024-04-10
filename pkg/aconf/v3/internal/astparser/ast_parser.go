package astparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type AstParser struct {
	file   *ast.File
	finder ValueFinder
}

func NewAstParser(filename string) (*AstParser, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}
	if f.Scope == nil {
		return nil, fmt.Errorf("file %s has nil scope", filename)
	}

	p := AstParser{
		file: f,
	}
	return &p, nil
}

// FindTagValuesInStructDecl ищет в исходном тексте файла, переданного парсеру объявление
// структуры с именем structname, затем проходит по всем полям объявленной структуры и ищет в них теги.
// Если тэг найден, то запускается finder и результат его работы складывается в map[string]any.
// Ключом в мапе будет значение, найденное finder, а значением - либо строка, если поле базовый тип, то есть
// int, string, slice, map, либо вложенная map[string]any, если поле является структурой.
// Гарантируется что в возвращаемой мапе конечные значения (листья дерева) всегда имеют тип string.
func (p *AstParser) FindTagValuesInStructDecl(structname string, finder ValueFinder) (map[string]any, error) {
	st, err := p.findStructTypeSpecByName(structname)
	if err != nil {
		return nil, err
	}
	p.finder = finder
	return p.recursiveFind(st)
}

func (p *AstParser) findStructTypeSpecByName(name string) (*ast.StructType, error) {
	var (
		obj  *ast.Object     // объект в скоупе файла
		spec *ast.TypeSpec   // объявление типа
		st   *ast.StructType // тип-структура, которую объявляет объект
		ok   bool
	)
	// ищем имя в скоупе файла
	obj, ok = p.file.Scope.Objects[name]
	if !ok {
		return nil, fmt.Errorf("name %s not found in file's scope", name)
	}
	if obj.Decl == nil {
		return nil, fmt.Errorf("object %s doe not declare anything", name)
	}
	// проверяем, что это объявление (declaration) типа (type X ...)
	spec, ok = obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("object %s declaration is not a type spec", name)
	}
	if spec.Type == nil {
		return nil, fmt.Errorf("object %s type is nil although is a type spec", name)
	}
	// проверяем, что объявленный тип - это структура (type X struct {})
	st, ok = spec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("object %s spec type is not a struct", name)
	}
	return st, nil
}

// recursiveFind обходит поля структуры рекурсивно, докапываясь до базовых типов (string, int, slice, map etc).
// Для полей, у которых есть теги, запускается ранее переданный паркеру finder и найденное значение
// складывается как ключ в мапу. Если у структуры нет тегов, то в качестве ключа добавляется ":N", где N -
// числа от 0 по порядку. Значением в мапе будет либо строкогово представление базового типа, либо
// вложенная мапа.
func (p *AstParser) recursiveFind(st *ast.StructType) (map[string]any, error) {
	var (
		result = make(map[string]any)
	)
	// создадим суррогатный ID для добавления в мапу
	// на случай, если тега у поля структуры нет, но он есть во вложенной структуре
	tagID := 0
	for _, f := range st.Fields.List {
		var (
			// если поле не базовый тип (string, int, slice), то сложим его
			// описание в отдельную мапу
			inner map[string]any
		)

		if innerStruct, ok := f.Type.(*ast.StructType); ok {
			// если внутри безымянная структура - поищем в ней
			inner, _ = p.recursiveFind(innerStruct)
		} else if innerStruct, err := p.findStructTypeSpecByName(fmt.Sprintf("%s", f.Type)); err == nil {
			// если полем является другой тип-структура, попробуем поискать его
			// объявление и найти теги в объявлении структуры
			inner, _ = p.recursiveFind(innerStruct)
		}

		var tagvalue string
		if f.Tag != nil {
			tagvalue = p.finder.GetValue(f.Tag.Value)
		} else {
			// если у поля не нашелся тег, положим в мапу искусственный ID.
			tagvalue = fmt.Sprintf(":%d", tagID)
			tagID++
		}

		if len(inner) > 0 {
			// мапа не пустая, значит у нас в поле лежит вложенная структура
			result[tagvalue] = inner
		} else {
			// поле - базовый тип
			result[tagvalue] = makeDefaultValue(f)
		}

	}
	return result, nil
}

// makeDefaultValue создаст строковое представление значения базового типа
func makeDefaultValue(f *ast.Field) string {
	var out string
	switch f.Type.(type) {
	case *ast.MapType:
		keyType := fmt.Sprintf("%s", f.Type.(*ast.MapType).Key)
		valueType := fmt.Sprintf("%s", f.Type.(*ast.MapType).Value)
		result := make([]string, 3)
		for i := 1; i < 4; i++ {
			key := fmt.Sprintf("%s%d", keyType, i)
			value := fmt.Sprintf("%s%d", valueType, i)
			result[i-1] = fmt.Sprintf("%s:%s", key, value)
		}
		out = strings.Join(result, ",")
	case *ast.ArrayType:
		elemType := fmt.Sprintf("%s", f.Type.(*ast.ArrayType).Elt)
		result := make([]string, 3)
		for i := 1; i < 4; i++ {
			result[i-1] = fmt.Sprintf("%s%d", elemType, i)
		}
		out = strings.Join(result, ",")
	default:
		out = fmt.Sprintf("%v", f.Type)
	}
	return out
}
