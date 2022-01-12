package lazy

import (
	"context"
	"sync"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/db"
	internal "github.com/ydb-platform/ydb-go-sdk/v3/internal/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/config"
)

type lazyTable struct {
	db      db.Connection
	options []config.Option
	client  table.Client
	m       sync.Mutex
}

func Table(db db.Connection, options []config.Option) *lazyTable {
	return &lazyTable{
		db:      db,
		options: options,
	}
}

func (t *lazyTable) CreateSession(ctx context.Context) (s table.ClosableSession, err error) {
	t.init(ctx)
	return t.client.CreateSession(ctx)
}

func (t *lazyTable) Do(ctx context.Context, op table.Operation, opts ...table.Option) (err error) {
	t.init(ctx)
	return t.client.Do(ctx, op, opts...)
}

func (t *lazyTable) DoTx(ctx context.Context, op table.TxOperation, opts ...table.Option) (err error) {
	t.init(ctx)
	return t.client.DoTx(ctx, op, opts...)
}

func (t *lazyTable) Close(ctx context.Context) error {
	t.m.Lock()
	defer t.m.Unlock()
	if t.client == nil {
		return nil
	}
	defer func() {
		t.client = nil
	}()
	return t.client.Close(ctx)
}

func (t *lazyTable) init(ctx context.Context) {
	t.m.Lock()
	if t.client == nil {
		t.client = internal.New(ctx, t.db, t.options...)
	}
	t.m.Unlock()
}
