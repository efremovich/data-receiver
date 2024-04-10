package tprepo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	postgres "github.com/efremovich/data-receiver/pkg/postgresdb"
)

// Поднимает посгрес в контейнере и накатывает на него миграции через goose. должен быть запущен docker и установлен goose.
// Прогоняет один и тот же тест как для реального репозитория в посгресе в контейнере, так и для моков репозиториев.
func TestTpRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgres.GetMockConn("../../../../migrations/package_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. обычно не включен докер, или неправильно указан путь к миграциям, или не установлен goose. %s", err.Error())
	}

	sqlTpRepo, err := NewTransportPackageRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	mockRepo := NewTpRepoMock()

	repos := []TransportPackageRepo{sqlTpRepo, mockRepo}

	for _, tpRepo := range repos {
		err = tpRepo.Ping(context.Background())
		if err != nil {
			t.Fatalf(err.Error())
		}

		name := uuid.NewString()
		origin := uuid.NewString()
		receiptURL := uuid.NewString()

		tpModel, err := tpRepo.Insert(ctx, name, receiptURL)
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, entity.TpStatusEnumNew, tpModel.Status)
		assert.Equal(t, name, tpModel.Name)
		assert.Equal(t, receiptURL, tpModel.ReceiptURL)

		tpModel, err = tpRepo.SelectByID(ctx, tpModel.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, entity.TpStatusEnumNew, tpModel.Status)
		assert.Equal(t, name, tpModel.Name)
		assert.Equal(t, receiptURL, tpModel.ReceiptURL)

		tpModel.Status = entity.TpStatusEnumSuccess
		tpModel.ErrorCode = "1111"
		tpModel.ErrorText = "text_error"
		tpModel.Origin = origin
		tr := true
		tpModel.IsReceipt = &tr

		err = tpRepo.UpdateExecOne(ctx, entity.TransportPackage{ID: 666})
		require.ErrorIs(t, sql.ErrNoRows, err)

		err = tpRepo.UpdateExecOne(ctx, *tpModel)
		if err != nil {
			t.Fatalf(err.Error())
		}

		_, err = tpRepo.SelectByName(ctx, "not_found")
		require.ErrorIs(t, err, sql.ErrNoRows)

		tpModel, err = tpRepo.SelectByName(ctx, name)
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, entity.TpStatusEnumSuccess, tpModel.Status)
		assert.Equal(t, "1111", tpModel.ErrorCode)
		assert.Equal(t, "text_error", tpModel.ErrorText)
		assert.Equal(t, origin, tpModel.Origin)
		assert.True(t, *tpModel.IsReceipt)

		tpModel.ErrorCode = ""
		tpModel.ErrorText = ""

		err = tpRepo.UpdateExecOne(ctx, *tpModel)
		if err != nil {
			t.Fatalf(err.Error())
		}

		tpModel, err = tpRepo.SelectByID(ctx, tpModel.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Empty(t, tpModel.ErrorCode)
		assert.Empty(t, tpModel.ErrorText)

		for _, event := range tpEventTypeEnumList {
			err = tpRepo.AddNewEvent(tpModel.ID, event, "")
			if err != nil {
				t.Fatalf(err.Error())
			}
		}

		res, err := tpRepo.SelectEvents(tpModel.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, len(tpEventTypeEnumList), len(res))
	}
}

func TestTransaction(t *testing.T) {
	conn, _, err := postgres.GetMockConn("../../../../migrations/package_receiver_db")
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx := context.Background()

	name := uuid.NewString()
	receiptURL := uuid.NewString()

	sqlTpRepo, err := NewTransportPackageRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tx, err := sqlTpRepo.BeginTX(context.TODO())
	if err != nil {
		t.Fatalf(err.Error())
	}

	tpModel, err := sqlTpRepo.WithTx(&tx).Insert(ctx, name, receiptURL)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tpWithoutTX, err := sqlTpRepo.SelectByID(ctx, tpModel.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(err.Error())
	}

	assert.Nil(t, tpWithoutTX)

	tpWithTX, err := sqlTpRepo.WithTx(&tx).SelectByID(ctx, tpModel.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, tpWithTX)

	tpWithTX.Status = entity.TpStatusEnumSuccess

	err = sqlTpRepo.WithTx(&tx).UpdateExecOne(ctx, *tpWithTX)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf(err.Error())
	}

	tpWithoutTX, err = sqlTpRepo.SelectByID(ctx, tpModel.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, tpWithoutTX)
	assert.Equal(t, name, tpWithoutTX.Name)
}

func TestSaveFileStructure(t *testing.T) {
	conn, _, err := postgres.GetMockConn("../../../../migrations/package_receiver_db")
	if err != nil {
		t.Fatalf(err.Error())
	}

	ctx := context.Background()

	name1 := uuid.NewString()
	name2 := uuid.NewString()
	receiptURL := uuid.NewString()

	sqlTpRepo, err := NewTransportPackageRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tp, err := sqlTpRepo.Insert(ctx, name1, receiptURL)
	if err != nil {
		t.Fatalf(err.Error())
	}

	r := "receipts.xml"
	d := "description.xml"
	lsID := "aaa55eb214b44e4eb18e163154af83a0"

	ex := []*entity.TpDirectory{{Name: ".", Files: map[string][]byte{r: nil, "txt.1": nil}}, {Name: lsID, Files: map[string][]byte{d: nil, "txt.2": nil}}}

	err = sqlTpRepo.SaveFileStructure(ctx, tp.ID, ex)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var res []struct {
		Dir string `db:"dir"`
		Doc string `db:"doc"`
	}

	err = conn.GetReadConnection().Select(&res, "SELECT d.name as dir, doc.name as doc FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id WHERE d.tp_id = $1", tp.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, res, 4)

	dirs, err := sqlTpRepo.SelectFileStructure(ctx, tp.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, dirs, 2)

	// При повторной вставке не должно измениться.
	err = sqlTpRepo.SaveFileStructure(ctx, tp.ID, ex)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = conn.GetReadConnection().Select(&res, "SELECT d.name as dir, doc.name as doc FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, res, 4)

	dirs, err = sqlTpRepo.SelectFileStructure(ctx, tp.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, dirs, 2)

	// Вставим то же самое для ещё одного ТП.
	tpSecond, err := sqlTpRepo.Insert(ctx, name2, receiptURL)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = sqlTpRepo.SaveFileStructure(ctx, tpSecond.ID, ex)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = conn.GetReadConnection().Select(&res, "SELECT d.name as dir, doc.name as doc FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, res, 8)

	dirs, err = sqlTpRepo.SelectFileStructure(ctx, tpSecond.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Len(t, dirs, 2)
}
