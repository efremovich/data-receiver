package usecases

import (
	"context"
	"fmt"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveData(ctx context.Context) aerror.AError {
	fmt.Println("implement me")
	return nil
}
