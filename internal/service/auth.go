package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	m "github.com/cs2-server/backend/internal/model"
	"golang.org/x/net/context"
)

type authStorage interface {
	GetProfileStatsByID(context.Context, string) (m.Stats, error)
}

type AuthService struct {
	storage authStorage
}

func NewAuthService(storage authStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (s *AuthService) GetProfile(ctx context.Context, apiKey string, ID string) (m.Profile, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", apiKey, ID)

	r, err := http.Get(url)
	if err != nil {
		return m.Profile{}, fmt.Errorf("GetProfile (1): %w", err)
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return m.Profile{}, fmt.Errorf("GetProfile (2): %w", err)
	}

	var player m.PlayerResponse
	err = json.Unmarshal(data, &player)
	if err != nil {
		return m.Profile{}, fmt.Errorf("GetProfile (3): %w", err)
	}

	if len(player.Response.Players) == 0 {
		return m.Profile{}, fmt.Errorf("GetProfile (4): %w", err)
	}

	p := player.Response.Players[0]

	stats, err := s.storage.GetProfileStatsByID(ctx, ID)
	if err != nil {
		return m.Profile{}, fmt.Errorf("GetProfile (5): %w", err)
	}

	return m.Profile{
		ID:           p.ID,
		Name:         p.Name,
		URL:          p.URL,
		Avatar:       p.Avatar,
		Kills:        stats.Kills,
		Deaths:       stats.Deaths,
		HeadshotRate: countHeadshotRate(stats.Kills, stats.Headshots),
	}, nil
}

func countHeadshotRate(kills int, headshots int) int {
	if kills <= 0 {
		return 0
	}

	hsRate := float64(headshots) / float64(kills) * 100

	return int(math.Round(hsRate))
}
