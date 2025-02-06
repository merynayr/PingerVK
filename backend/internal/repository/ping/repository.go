package ping

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/merynayr/PingerVK/backend/internal/model"
	"github.com/merynayr/PingerVK/backend/internal/repository"
	repoModel "github.com/merynayr/PingerVK/backend/internal/repository/model"
	"github.com/merynayr/PingerVK/pkg/client/db"
	"github.com/merynayr/PingerVK/pkg/logger"

	sq "github.com/Masterminds/squirrel"
)

// Таблицы
const (
	containersTable = "containers"
	pingsTable      = "pings"

	// Колонки containers
	containerIDColumn = "id_container"
	NameColumn        = "name"

	// Колонки pings
	containerFKColumn  = "id_container"
	IPColumn           = "ip"
	statusColumn       = "status"
	responseTimeColumn = "response_time"
	lastSuccessColumn  = "last_success"
)

// Репозиторий
type repo struct {
	db db.Client
}

// NewRepository возвращает объект репозитория
func NewRepository(db db.Client) repository.PingRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, ping *model.Pings) (int64, error) {
	op := "Ping.Create"

	query, args, err := sq.Insert(containersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(containerIDColumn, NameColumn).
		Values(ping.ID, ping.Name).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		logger.Error("%s: failed to build container insert query: %v", op, err)
		return 0, err
	}

	q := db.Query{
		Name:     "ping_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		logger.Error("%s: failed to insert container: %v", q.Name, err)
		return 0, err
	}

	query, args, err = sq.Insert(pingsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(containerFKColumn, IPColumn, statusColumn, responseTimeColumn, lastSuccessColumn).
		Values(
			ping.ID,
			ping.IPAddress,
			ping.Status,
			ping.ResponseTime/1000000,
			ping.LastSuccess,
		).
		ToSql()

	if err != nil {
		logger.Error("%s: failed to build ping insert query: %v", op, err)
		return 0, err
	}

	q = db.Query{
		Name:     "ping_repository.CreatePing",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error("%s: failed to insert ping: %v", q.Name, err)
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context) ([]*model.GetPings, error) {
	op := "Ping.Get"

	query, args, err := sq.Select(
		"c.id_container", "c.name",
		"p.ip", "p.status", "p.response_time", "p.last_success",
	).
		From("pings p").
		Join("containers c ON p.id_container = c.id_container").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		logger.Error("%s: failed to create query: %v", op, err)
		return nil, err
	}

	// Выполняем запрос
	q := db.Query{
		Name:     "ping_repository.Get",
		QueryRaw: query,
	}

	var result []repoModel.Pings
	err = r.db.DB().ScanAllContext(ctx, &result, q, args...)
	if err != nil {
		logger.Error("%s: failed to fetch ping data: %v", op, err)
		return nil, err
	}

	// Преобразуем результаты в []*model.Pings
	var pings []*model.GetPings
	for _, res := range result {
		pings = append(pings, &model.GetPings{
			ID:           res.ID,
			Name:         res.Name,
			IPAddress:    res.IPAddress,
			Status:       res.Status,
			ResponseTime: res.ResponseTime,
			LastSuccess:  res.LastSuccess.UTC().In(time.Local).Format("2006-01-02 15:04:05"),
		})
	}

	return pings, nil
}

func (r *repo) ContainerExists(ctx context.Context, id string) (bool, error) {
	query, args, err := sq.Select("1").
		From(containersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{containerIDColumn: id}).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build chat check query: %w", err)
	}

	q := db.Query{
		Name:     "ping_repository.ContainerExists",
		QueryRaw: query,
	}

	var containerID string
	err = r.db.DB().ScanOneContext(ctx, &containerID, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return true, nil
}

func (r *repo) UpdateInfo(ctx context.Context, ping *model.Pings) error {
	op := "Ping.UpdateInfo"

	lastSuccessExpr := sq.Expr(lastSuccessColumn)
	if ping.Status {
		lastSuccessExpr = sq.Expr("?::timestamp", ping.LastSuccess.UTC().In(time.Local).Format("2006-01-02 15:04:05"))
	}

	queryUpdate, argsUpdate, err := sq.Update(pingsTable).
		Set(IPColumn, ping.IPAddress).
		Set(statusColumn, ping.Status).
		Set(responseTimeColumn, ping.ResponseTime/1000000).
		Set(lastSuccessColumn, lastSuccessExpr).
		Where(sq.Eq{containerFKColumn: ping.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		logger.Error("%s: failed to build update query: %v", op, err)
		return err
	}

	q := db.Query{
		Name:     "ping_repository.UpdatePing",
		QueryRaw: queryUpdate,
	}

	_, err = r.db.DB().ExecContext(ctx, q, argsUpdate...)
	if err != nil {
		logger.Error("%s: failed to update ping: %v", q.Name, err)
		return err
	}

	return nil
}
