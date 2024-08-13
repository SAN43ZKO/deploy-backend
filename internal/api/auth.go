package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	steamauth "github.com/TeddiO/GoSteamAuth/src"
	"github.com/cs2-server/backend/config"
	m "github.com/cs2-server/backend/internal/model"
	"github.com/cs2-server/backend/internal/render"
	"github.com/sirupsen/logrus"
)

const (
	ErrMethodNotAllowed = "method not allowed"
	ErrInvalidAuth      = "invalid auth"
	ErrParamNotSet      = "param is not set"
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

// @Summary Redirects client to Steam authentication page
// @Tags auth
// @Accept json
// @Produce json
// @Success 302 {object} nil
// @Failure 405 {object} render.Err
// @Router /api/auth/login [get]
func (a *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		render.Error(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)

		return
	}

	query := fmt.Sprintf("http://%s:%s/api/auth/process", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
	steamauth.RedirectClient(w, r, steamauth.BuildQueryString(query))
}

// @Summary Processes Steam authentication response and generates JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} m.JWT
// @Failure 400 {object} render.Err
// @Failure 405 {object} render.Err
// @Failure 500 {object} render.Err
// @Router /api/auth/process [get]
func (a *AuthAPI) ProcessLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		render.Error(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)

		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		a.logger.Errorln(err)
		render.Error(w, http.StatusInternalServerError, err.Error())

	}

	queryMap := steamauth.ValuesToMap(query)

	steamID, isValid, err := steamauth.ValidateResponse(queryMap)
	if err != nil {
		a.logger.Errorln(err)
		render.Error(w, http.StatusInternalServerError, err.Error())

		return
	}

	if !isValid {
		a.logger.Errorln(ErrInvalidAuth)
		render.Error(w, http.StatusInternalServerError, ErrInvalidAuth)

		return
	}

	tokens, err := a.jwt.GenerateTokens(steamID)
	if err != nil {
		a.logger.Errorln(err)
		render.Error(w, http.StatusInternalServerError, ErrInvalidAuth)
	}

	render.JSON(w, http.StatusOK, tokens)
}

// @Summary Refreshes JWT tokens
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Param id formData string true "User ID"
// @Success 200 {object} m.JWT
// @Failure 400 {object} render.Err
// @Failure 401 {object} render.Err
// @Failure 405 {object} render.Err
// @Failure 500 {object} render.Err
// @Router /api/auth/refresh [post]
func (a *AuthAPI) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.logger.Errorln(ErrMethodNotAllowed)
		render.Error(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)

		return
	}

	id := r.FormValue("id")

	if id == "" {
		a.logger.Errorln(ErrParamNotSet)
		render.Error(w, http.StatusBadRequest, ErrParamNotSet)

		return
	}

	tokens, err := a.jwt.GenerateTokens(id)
	if err != nil {
		a.logger.Errorln(err)
		render.Error(w, http.StatusInternalServerError, err.Error())
	}

	render.JSON(w, http.StatusOK, tokens)
}

// @Summary Retrieves user profile
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} m.Profile
// @Failure 400 {object} render.Err
// @Failure 401 {object} render.Err
// @Failure 500 {object} render.Err
// @Router /api/profile/{id} [get]
func (a *AuthAPI) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.logger.Errorln(ErrMethodNotAllowed)
		render.Error(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)

		return
	}

	id := r.PathValue("id")

	profile, err := a.service.GetProfile(r.Context(), a.cfg.Steam.APIKey, id)
	if err != nil {
		a.logger.Errorln(err)
		render.Error(w, http.StatusInternalServerError, err.Error())

		return
	}

	render.JSON(w, http.StatusOK, profile)
}
