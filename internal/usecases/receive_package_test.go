package usecases

import (
	"context"
	"path"
	"testing"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestBuildFileStructure(t *testing.T) {
	type testCase struct {
		name string
		in   map[string][]byte
		out  []entity.TpDirectory
	}

	testCases := []testCase{
		{name: "1", in: nil, out: []entity.TpDirectory{}},
		{name: "2", in: map[string][]byte{
			"receipt.xml": nil,
			"1.txt":       nil,
			"2.txt":       nil,
			"dir1/3.txt":  nil,
		}, out: []entity.TpDirectory{
			{Name: ".", Files: map[string][]byte{"receipt.xml": nil, "1.txt": nil, "2.txt": nil}},
			{Name: "dir1", Files: map[string][]byte{"3.txt": nil}},
		}},
		{name: "3", in: map[string][]byte{
			"dir1/receipt.xml": nil,
			"dir1/1.txt":       nil,
			"dir2/2.txt":       nil,
			"dir3/3.txt":       nil,
		}, out: []entity.TpDirectory{
			{Name: "dir1", Files: map[string][]byte{"receipt.xml": nil, "1.txt": nil}},
			{Name: "dir2", Files: map[string][]byte{"2.txt": nil}},
			{Name: "dir3", Files: map[string][]byte{"3.txt": nil}},
		}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			res := buildFileStructure(context.Background(), test.in)

			assert.Len(t, test.out, len(res))

			for i, v := range res {
				assert.Equal(t, test.out[i], *v)
			}
		})
	}
}

func TestValidateFilesStructure(t *testing.T) {
	type testCase struct {
		name  string
		in    []*entity.TpDirectory
		valid bool
	}

	r := "receipts.xml"
	d := "description.xml"
	lsID := "aaa55eb214b44e4eb18e163154af83a0"

	filaName1 := "d888c4d7563b4a218d1da80843da4b8a.p7s"
	filaName2 := "a0587509136b4a00b0dc59b79a144398.bin"

	testCases := []testCase{
		{name: "пусто", in: nil, valid: false},
		{name: "квитанция_и_лс", in: []*entity.TpDirectory{
			{Name: ".", Files: map[string][]byte{r: nil}},
			{Name: lsID, Files: map[string][]byte{d: nil}},
		}, valid: false},
		{name: "в_корне_нет_квитанции", in: []*entity.TpDirectory{
			{Name: ".", Files: map[string][]byte{filaName1: nil}},
		}, valid: false},
		{name: "в_лc_нет_описания", in: []*entity.TpDirectory{
			{Name: lsID, Files: map[string][]byte{filaName1: nil}},
		}, valid: false},
		{name: "валидная_квитанция", in: []*entity.TpDirectory{
			{Name: ".", Files: map[string][]byte{r: nil, filaName1: nil, filaName2: nil}},
		}, valid: true},
		{name: "валидный_лс", in: []*entity.TpDirectory{
			{Name: lsID, Files: map[string][]byte{d: nil, filaName1: nil, filaName2: nil}},
		}, valid: true},
		{name: "невалидное_название_файла", in: []*entity.TpDirectory{
			{Name: lsID, Files: map[string][]byte{d: nil, filaName1: nil, "123.txt": nil}},
		}, valid: false},
		{name: "двойная_вложенность", in: []*entity.TpDirectory{
			{Name: path.Join("aaa55eb214b44", "e4eb18e163154af83a0"), Files: map[string][]byte{d: nil, filaName1: nil, filaName2: nil}},
		}, valid: false},
		{name: "невалидное_имя_директории", in: []*entity.TpDirectory{
			{Name: "dir_1", Files: map[string][]byte{d: nil, filaName1: nil, filaName2: nil}},
		}, valid: false},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := validateFilesStructure(context.Background(), test.in)
			assert.Equal(t, test.valid, err == nil)
		})
	}
}
