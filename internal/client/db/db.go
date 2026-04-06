package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Клиент для работы с бд
type Client interface {
	DB() DB
	Close() error
}

// Query обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться еще где-то, например, в тестах
type Query struct {
	Name     string // имя метода к примеру
	QueryRaw string // squrrel выплюнет сырой запрос
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Handler func(ctx context.Context) error

type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// Интерфейс для работы с именованными запросами с помощью тегов с структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// Интерфейс для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}
