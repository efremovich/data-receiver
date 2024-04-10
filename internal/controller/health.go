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
		err := gw.packageReceiver.PingDB(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге БД приёмника", err, nil, false)
			return status.Errorf(codes.Internal, "ошибка при пинге БД приёмника: %s", err.Error())
		}

		return nil
	})

	g.Go(func() error {
		err := gw.packageReceiver.PingNATS(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге НАТС", err, nil, false)
			return status.Errorf(codes.Internal, "ошибка при пинге НАТС: %s", err.Error())
		}

		return nil
	})

	g.Go(func() error {
		err := gw.packageReceiver.PingStorage(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге стораджа", err, nil, false)
			return status.Errorf(codes.Internal, "ошибка при пинге стораджа: %s", err.Error())
		}

		return nil
	})

	g.Go(func() error {
		err := gw.packageReceiver.PingOperator(ctx)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка при пинге апи опертора", err, nil, false)
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
