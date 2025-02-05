package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/merynayr/PingerVK/pkg/client/db"
	"github.com/pkg/errors"
)

// pgClient структура клиента postgres
type pgClient struct {
	masterDBC db.DB
}

// New возвращает новый объект клиента для работы с postgres
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: NewDB(dbc),
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
