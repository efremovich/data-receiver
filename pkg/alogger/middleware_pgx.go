package alogger

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	pgxQueryStartedAtCtxValueKey = "_pgxQueryStartedAt"
	pgxTraceQueryData            = "_pgxTraceQueryData"
)

type PGXQueryTracer struct {
	SkipFirstArg bool // При использовании db/sql или sqlx - ставить true. Если использовать чистый pgx, то оставить false
}

func (p *PGXQueryTracer) trimUnnecessaryArgs(args []any) []any {
	if p.SkipFirstArg && len(args) > 1 {
		return args[1:]
	}

	return args
}

func (p *PGXQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log := NewALogger(ctx, nil)
	log.Debugf("SQL START").
		SetAttr("sql", data.SQL).SetAttr("args", p.trimUnnecessaryArgs(data.Args))
	_ = log.Flush()

	ctx = context.WithValue(ctx, pgxTraceQueryData, data)
	ctx = context.WithValue(ctx, pgxQueryStartedAtCtxValueKey, time.Now())

	return ctx
}

func (p *PGXQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	startedAt, ok := ctx.Value(pgxQueryStartedAtCtxValueKey).(time.Time)
	if !ok {
		return
	}

	dur := time.Since(startedAt)

	queryData, ok := ctx.Value(pgxTraceQueryData).(pgx.TraceQueryStartData)
	if !ok {
		return
	}

	log := NewALogger(ctx, nil)
	log.Debugf("SQL END").
		SetAttr("sql", queryData.SQL).SetAttr("args", p.trimUnnecessaryArgs(queryData.Args)).SetAttr("time", dur)

	_ = log.Flush()
}
