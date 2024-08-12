package storage

import (
	"context"
	"fmt"

	m "github.com/cs2-server/backend/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthStorage struct {
	db *pgxpool.Pool
}

func NewAuthStorage(db *pgxpool.Pool) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

func (s *AuthStorage) GetProfileStatsByID(ctx context.Context, ID string) (m.Stats, error) {
	query := `
        SELECT kills, deaths, headshots
        FROM player_stats
        WHERE steam_id = $1
    `

	var stats m.Stats
	if err := s.db.QueryRow(ctx, query, ID).Scan(&stats.Kills, &stats.Deaths, &stats.Headshots); err != nil {
		return m.Stats{}, fmt.Errorf("GetProfileStats: %w", err)
	}

	return stats, nil
}
