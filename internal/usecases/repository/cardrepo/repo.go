package cardrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var _ CardRepo = (*cardRepo)(nil)
type CardRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Card, error)
	SelectByVendorID(ctx context.Context, vendorID string) (*entity.Card, error)
	SelectByVendorCode(ctx context.Context, vendorCode string) (*entity.Card, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Card, error)

	Insert(ctx context.Context, in entity.Card) error
	UpdateExecOne(ctx context.Context, in entity.Card) error

	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CardRepo
}

type cardRepoImpl struct{

}
