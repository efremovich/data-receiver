package tprepo

import (
	"context"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"

	postgres "github.com/efremovich/data-receiver/pkg/postgresdb"
)

type TransportPackageRepo interface {
	SelectByID(ctx context.Context, tpID int64) (*entity.TransportPackage, error)
	SelectByName(ctx context.Context, name string) (*entity.TransportPackage, error)

	SelectByDocument(ctx context.Context, doc string) (*entity.TransportPackage, error)

	Insert(ctx context.Context, name string, receiptURL string) (*entity.TransportPackage, error)
	UpdateExecOne(ctx context.Context, tp entity.TransportPackage) error

	// Сохраняет в БД информацию о директориях и документах в ТП.
	SaveFileStructure(ctx context.Context, tpID int64, fileStructure []*entity.TpDirectory) error
	SelectFileStructure(ctx context.Context, tpID int64) ([]*entity.TpDirectory, error)

	Ping(ctx context.Context) error

	AddNewEvent(tpID int64, event entity.TpEventTypeEnum, desc string) error
	SelectEvents(tpID int64) ([]entity.TpEvent, error)

	BeginTX(ctx context.Context) (postgres.Transaction, error)
	WithTx(*postgres.Transaction) TransportPackageRepo
}

type transportPackageRepoImpl struct {
	statusesMap map[entity.TpStatusEnum]int
	eventsMap   map[entity.TpEventTypeEnum]int
	db          *postgres.DBConnection
	tx          *postgres.Transaction
}

func NewTransportPackageRepo(_ context.Context, db *postgres.DBConnection) (TransportPackageRepo, error) {
	statuses := []struct {
		ID   int    `db:"id"`
		Desc string `db:"tp_status_desc"`
	}{}

	err := db.GetReadConnection().Select(&statuses, "SELECT id, tp_status_desc FROM tp_status_enum")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке статусов ТП из бд  - %s", err.Error())
	}

	statusesMap := make(map[entity.TpStatusEnum]int)

	for _, status := range statuses {
		var found bool

		for _, codeStatus := range tpStatusEnumList {
			if codeStatus == entity.TpStatusEnum(status.Desc) {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("enum статусов ТП в коде не соответствует таблице tp_status_enum. не найден %s", status.Desc)
		}

		statusesMap[entity.TpStatusEnum(status.Desc)] = status.ID
	}

	if len(tpStatusEnumList) != len(statusesMap) {
		return nil, fmt.Errorf("enum статусов ТП в коде не соответствует таблице tp_status_enum. длина enum в коде - %d, в БД - %d", len(tpStatusEnumList), len(statusesMap))
	}

	events := []struct {
		ID   int    `db:"id"`
		Desc string `db:"tp_event_desc"`
	}{}

	err = db.GetReadConnection().Select(&events, "SELECT id, tp_event_desc FROM tp_event_enum")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выборке евентов ТП из бд  - %s", err.Error())
	}

	eventsMap := make(map[entity.TpEventTypeEnum]int)

	for _, event := range events {
		var found bool

		for _, codeEvent := range tpEventTypeEnumList {
			if codeEvent == entity.TpEventTypeEnum(event.Desc) {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("enum евентов ТП в коде не соответствует таблице tp_event_enum. не найден %s", event.Desc)
		}

		eventsMap[entity.TpEventTypeEnum(event.Desc)] = event.ID
	}

	if len(tpEventTypeEnumList) != len(eventsMap) {
		return nil, fmt.Errorf("enum евентов ТП в коде не соответствует таблице tp_event_enum")
	}

	return &transportPackageRepoImpl{db: db, statusesMap: statusesMap, eventsMap: eventsMap}, nil
}

const (
	queryDocflowFullSelect = `
	SELECT 
		tp.id,
		tp.name,
		tp.is_receipt,
		tp.sender_operator_code,
		tp.receipt_url,
		tp.created_at,
		status.tp_status_desc,
		tperr.error_text,
		tperr.error_code
	FROM transport_package tp 
		LEFT JOIN tp_error tperr ON tperr.tp_id = tp.id
		LEFT JOIN tp_status_enum status ON tp.tp_status_id = status.id `
)

func (repo *transportPackageRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.TransportPackage, error) {
	var result transportPackageDB

	query := queryDocflowFullSelect + `WHERE tp.id = $1`

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}

	return result.ConvertToEntityTransportPackage(ctx), nil
}

func (repo *transportPackageRepoImpl) SelectByName(ctx context.Context, name string) (*entity.TransportPackage, error) {
	var result transportPackageDB

	query := queryDocflowFullSelect + `WHERE tp.name = $1`

	err := repo.getReadConnection().Get(&result, query, name)
	if err != nil {
		return nil, err
	}

	return result.ConvertToEntityTransportPackage(ctx), nil
}

func (repo *transportPackageRepoImpl) SelectByDocument(ctx context.Context, doc string) (*entity.TransportPackage, error) {
	var tpID repository.IDWrapper

	err := repo.getReadConnection().Get(&tpID, "SELECT d.tp_id as id FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id WHERE doc.name = $1", doc)
	if err != nil {
		return nil, err
	}

	return repo.SelectByID(ctx, tpID.ID.Int64)
}

func (repo *transportPackageRepoImpl) Insert(_ context.Context, name string, receiptURL string) (*entity.TransportPackage, error) {
	tp := entity.TransportPackage{
		Name:       name,
		ReceiptURL: receiptURL,
		Status:     entity.TpStatusEnumNew,
	}

	query := `
		INSERT INTO transport_package (name, receipt_url, tp_status_id) 
		VALUES ($1, $2, $3) RETURNING id`

	tpIDdWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&tpIDdWrap, query, tp.Name, tp.ReceiptURL, repo.statusesMap[tp.Status])
	if err != nil {
		return nil, err
	}

	tp.ID = tpIDdWrap.ID.Int64

	return &tp, nil
}

func (repo *transportPackageRepoImpl) UpdateExecOne(ctx context.Context, tp entity.TransportPackage) error {
	dbModel := convertToDBTransportPackage(ctx, tp, repo.statusesMap)

	if dbModel.ErrorCode.Valid {
		_, err := repo.getWriteConnection().ExecOne(
			`INSERT INTO tp_error (tp_id, error_text, error_code) VALUES ($1, $2, $3) 
			ON CONFLICT (tp_id) DO UPDATE SET error_text = $2, error_code = $3`,
			dbModel.ID, dbModel.ErrorText, dbModel.ErrorCode)
		if err != nil {
			return err
		}
	} else {
		_, err := repo.getWriteConnection().Exec("DELETE FROM tp_error WHERE tp_id = $1", dbModel.ID)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE transport_package SET
			tp_status_id = $2,
			sender_operator_code = $3,
			is_receipt = $4,
			updated_at = current_timestamp
		WHERE id = $1`

	_, err := repo.getWriteConnection().ExecOne(query, dbModel.ID, dbModel.StatusID, tp.Origin, tp.IsReceipt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *transportPackageRepoImpl) SaveFileStructure(_ context.Context, tpID int64, fileStructure []*entity.TpDirectory) error {
	queryDirectory := `INSERT INTO tp_directory (tp_id, name) VALUES ($1, $2) ON CONFLICT (tp_id, name) DO UPDATE SET created_at = now() RETURNING id`
	queryDocument := `INSERT INTO tp_document (directory_id, name) VALUES ($1, $2) ON CONFLICT (directory_id, name) DO NOTHING`
	directoryID := repository.IDWrapper{}

	for i := range fileStructure {
		err := repo.getWriteConnection().QueryAndScan(&directoryID, queryDirectory, tpID, fileStructure[i].Name)
		if err != nil {
			return err
		}

		for filename := range fileStructure[i].Files {
			_, err := repo.getWriteConnection().Exec(queryDocument, directoryID.ID, filename)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *transportPackageRepoImpl) SelectFileStructure(_ context.Context, tpID int64) ([]*entity.TpDirectory, error) {
	var dirs []struct {
		Dir string `db:"dir"`
		Doc string `db:"doc"`
	}

	err := repo.getReadConnection().Select(&dirs, "SELECT d.name as dir, doc.name as doc FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id WHERE d.tp_id = $1", tpID)
	if err != nil {
		return nil, err
	}

	res := make([]*entity.TpDirectory, 0, len(dirs))

	for _, d1 := range dirs {
		var found *entity.TpDirectory

		for _, d2 := range res {
			if d2.Name == d1.Dir {
				found = d2
				break
			}
		}

		if found == nil {
			found = &entity.TpDirectory{
				Name:  d1.Dir,
				Files: make(map[string][]byte),
			}

			res = append(res, found)
		}

		found.Files[d1.Doc] = nil
	}

	return res, nil
}

func (repo *transportPackageRepoImpl) AddNewEvent(tpID int64, event entity.TpEventTypeEnum, desc string) error {
	query := `INSERT INTO tp_event (tp_id, event_type_id, description) VALUES ($1, $2, NULLIF($3, ''));`

	_, err := repo.getWriteConnection().ExecOne(query, tpID, repo.eventsMap[event], desc)
	if err != nil {
		return err
	}

	return nil
}

func (repo *transportPackageRepoImpl) SelectEvents(tpID int64) ([]entity.TpEvent, error) {
	query := ` SELECT tp_e.tp_id, tp_e.created_at, COALESCE(tp_e.description, '') AS description, ev_enum.tp_event_desc AS event_type
		FROM tp_event tp_e INNER JOIN tp_event_enum ev_enum ON tp_e.event_type_id = ev_enum.id
		WHERE tp_e.tp_id = $1
		ORDER BY tp_e.created_at ASC;`

	var res []entity.TpEvent

	err := repo.getReadConnection().Select(&res, query, tpID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *transportPackageRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *transportPackageRepoImpl) BeginTX(ctx context.Context) (postgres.Transaction, error) {
	return repo.db.GetWriteConnection().BeginTX(ctx)
}

func (repo *transportPackageRepoImpl) WithTx(tx *postgres.Transaction) TransportPackageRepo {
	return &transportPackageRepoImpl{db: repo.db, tx: tx, statusesMap: repo.statusesMap, eventsMap: repo.eventsMap}
}

func (repo *transportPackageRepoImpl) getReadConnection() postgres.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *transportPackageRepoImpl) getWriteConnection() postgres.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
