package db

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/omerkaya1/watcher/internal/config"
	"github.com/omerkaya1/watcher/internal/errors"
	"github.com/omerkaya1/watcher/internal/models"
	"time"
)

// MainEventStorage object holds everything related to the DB interactions
type MainEventStorage struct {
	db *sqlx.DB
}

// NewMainEventStorage returns new MainEventStorage object to the callee
func NewMainEventStorage(cfg config.DBConf) (*MainEventStorage, error) {
	if cfg.Name == "" || cfg.User == "" || cfg.SSLMode == "" || cfg.Password == "" {
		return nil, errors.ErrBadDBConfiguration
	}
	dsn := fmt.Sprintf("host=%s port=%s password=%s user=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Password, cfg.User, cfg.Name, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &MainEventStorage{db: db}, nil
}

func (mes *MainEventStorage) GetUpcomingEvents(ctx context.Context) ([]models.Event, error) {
	eventList := make([]models.Event, 0)
	e := models.Event{}
	query := "select * from events where start_time::date=current_date or start_time::date=current_date+interval '1 day'"
	rows, err := mes.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if err := rows.StructScan(&e); err != nil {
				return eventList, err
			}
			eventList = append(eventList, e)
		}
	}
	return eventList, nil
}

func (mes *MainEventStorage) GetEventsForSpecifiedDate(ctx context.Context, t time.Time) ([]models.Event, error) {
	return nil, nil
}
