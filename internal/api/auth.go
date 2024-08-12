package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	steamauth "github.com/TeddiO/GoSteamAuth/src"
	"github.com/cs2-server/backend/config"
	m "github.com/cs2-server/backend/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	ErrMethodNotAllowed = "method not allowed"
)

type jwtGenerator interface {
	GenerateTokens(string) (m.JWT, error)
}

type authService interface {
	GetProfile(context.Context, string, string) (m.Profile, error)
}

type AuthAPI struct {
	cfg     *config.Config
	logger  *logrus.Logger
	jwt     jwtGenerator
	service authService
}

func NewAuthAPI(cfg *config.Config, logger *logrus.Logger, jwt jwtGenerator, service authService) *AuthAPI {
	return &AuthAPI{
		cfg:     cfg,
		logger:  logger,
		jwt:     jwt,
		service: service,
	}
}

// @Summary Redirects the client to the Steam authentication page.
// @Description Redirects the client to the Steam authentication page using the Steam API.
// @Tags auth
// @Accept json
// @Produce json
// @Success 302 {object} nil
// @Failure 405 {object} nil
// @Router /auth/login [get]
func (a *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		JSON(w, http.StatusMethodNotAllowed, nil)

		return
	}

	query := fmt.Sprintf("http://%s:%s/auth/process", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
	steamauth.RedirectClient(w, r, steamauth.BuildQueryString(query))
}

// @Summary Processes the Steam authentication response and generates JWT tokens.
// @Description Validates the Steam authentication response and generates JWT tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} m.JWT
// @Failure 400 {object} nil
// @Failure 405 {object} nil
// @Failure 500 {object} nil
// @Router /auth/process [get]
func (a *AuthAPI) ProcessLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		JSON(w, http.StatusMethodNotAllowed, nil)

		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		a.logger.Errorln(err)
		JSON(w, http.StatusInternalServerError, nil)

	}

	queryMap := steamauth.ValuesToMap(query)

	steamID, isValid, err := steamauth.ValidateResponse(queryMap)
	if err != nil {
		a.logger.Errorln(err)
		JSON(w, http.StatusInternalServerError, nil)

		return
	}

	if !isValid {
		a.logger.Errorln("invalid auth")
		JSON(w, http.StatusInternalServerError, nil)

		return
	}

	tokens, err := a.jwt.GenerateTokens(steamID)
	if err != nil {
		a.logger.Errorln(err)
		JSON(w, http.StatusInternalServerError, nil)
	}

	JSON(w, http.StatusOK, tokens)
}

// @Summary Refreshes the JWT tokens.
// @Description Refreshes the JWT tokens based on the provided steam_id. Requires a valid JWT refresh token.
// @Tags auth
// @Security BearerAuth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param steam_id formData string true "Steam ID"
// @Success 200 {object} m.JWT
// @Failure 400 {object} nil
// @Failure 401 {object} nil
// @Failure 405 {object} nil
// @Failure 500 {object} nil
// @Router /auth/refresh [post]
func (a *AuthAPI) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.logger.Errorln(ErrMethodNotAllowed)
		JSON(w, http.StatusMethodNotAllowed, nil)

		return
	}

	steamID := r.FormValue("steam_id")

	if steamID == "" {
		a.logger.Errorln("steam_id param is not set")
		JSON(w, http.StatusBadRequest, nil)

		return
	}

	tokens, err := a.jwt.GenerateTokens(steamID)
	if err != nil {
		a.logger.Errorln(err)
		JSON(w, http.StatusInternalServerError, nil)
	}

	JSON(w, http.StatusOK, tokens)
}

// @Summary Retrieves the Steam user profile.
// @Description Fetches the Steam user profile based on the provided steam_id. Requires a valid JWT token.
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param steam_id query string true "Steam ID"
// @Success 200 {object} m.Profile
// @Failure 400 {object} nil
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /profile [get]
func (a *AuthAPI) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		JSON(w, http.StatusMethodNotAllowed, nil)

		return
	}

	steamID := r.FormValue("steam_id")

	if steamID == "" {
		a.logger.Errorln("steam_id param is not set")
		JSON(w, http.StatusBadRequest, nil)

		return
	}

	profile, err := a.service.GetProfile(r.Context(), a.cfg.Steam.APIKey, steamID)
	if err != nil {
		a.logger.Errorln(err)
		JSON(w, http.StatusInternalServerError, nil)

		return
	}

	JSON(w, http.StatusOK, profile)
}
