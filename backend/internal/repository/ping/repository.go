package ping

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
	"github.com/merynayr/PingerVK/backend/internal/repository"
	"github.com/merynayr/PingerVK/pkg/client/db"
	"github.com/merynayr/PingerVK/pkg/logger"

	sq "github.com/Masterminds/squirrel"
)

// Алиасы для базы данных
const (
	tableName = "pings"

	IDColumn           = "id"
	IPColumn           = "ip"
	statusColumn       = "status"
	responseTimeColumn = "response_time"
	lastSuccessColumn  = "last_success"
)

// Структура репо с клиентом базы данных (интерфейсом)
type repo struct {
	db db.Client
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db db.Client) repository.PingRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, ping *model.Pings) (int64, error) {
	op := "Ping.Create"

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(IPColumn, statusColumn, responseTimeColumn, lastSuccessColumn).
		Values(ping.IP, ping.Status, ping.ResponseTime, ping.LastSuccess).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		logger.Error("%s: failed to create builder: %v", op, err)
		return 0, err
	}

	q := db.Query{
		Name:     "ping_repository.Create",
		QueryRaw: query,
	}

	var ID int64
	err = r.db.DB().ScanOneContext(ctx, &ID, q, args...)
	if err != nil {
		logger.Error("%s: failed to insert ping information: %v", op, err)
		return 0, err
	}

	return ID, nil
}

func (r *repo) Get(ctx context.Context) ([]*model.Pings, error) {
	op := "Ping.Get"

	// Создаем запрос для получения всех пингов
	query, args, err := sq.Select(IPColumn, statusColumn, responseTimeColumn, lastSuccessColumn).
		From(tableName).
		ToSql()

	if err != nil {
		logger.Error("%s: failed to create builder: %v", op, err)
		return nil, err
	}

	// Выполняем запрос
	q := db.Query{
		Name:     "ping_repository.Get",
		QueryRaw: query,
	}

	var pings []*model.Pings
	err = r.db.DB().ScanAllContext(ctx, &pings, q, args...)
	if err != nil {
		logger.Error("%s: failed to fetch ping information: %v", op, err)
		return nil, err
	}

	return pings, nil
}
