package controller

import (
	"context"

	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/gogo/status"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

// хелз-чек для кубера. кубер смотрит только на статус-код ответа. если он не 200 - начинает рестартить поду.
// соответсвенно логика такая: если без чего-то сервис не может работать полностью (например без бд), то пусть падает
// если например недоступно api оператора - то частично функциональность сохраняется, поэтому 200.
func (gw *grpcGatewayServerImpl) CheckHealth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	g := errgroup.Group{}

	g.Go(func() error {
		err := gw.core.PingDB(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге БД приёмника %s", err.Error())
			return status.Errorf(codes.Internal, "ошибка при пинге БД приёмника: %s", err.Error())
		}

		return nil
	})

	g.Go(func() error {
		err := gw.core.PingNATS(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге НАТС %s", err.Error())
			return status.Errorf(codes.Internal, "ошибка при пинге НАТС: %s", err.Error())
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
