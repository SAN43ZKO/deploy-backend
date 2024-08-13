package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	m "github.com/cs2-server/backend/internal/model"
	"github.com/cs2-server/backend/internal/render"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenExpired = errors.New("token has expired")
)

type JWT struct {
	key string
}

func New(key string) *JWT {
	return &JWT{
		key: key,
	}
}

func (t *JWT) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromHeader(r)
		if err != nil {
			logrus.Errorln("JWT (1):", err)
			render.Error(w, http.StatusUnauthorized, err.Error())

			return
		}

		if err := t.verifyToken(tokenString); err != nil {
			logrus.Errorln("JWT (2): ", err)

			if errors.Is(err, ErrTokenExpired) {
				render.Error(w, http.StatusUnauthorized, err.Error(), render.ExpiredToken)

				return
			}

			render.Error(w, http.StatusUnauthorized, err.Error())

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (t *JWT) GenerateTokens(id string) (m.JWT, error) {
	var (
		accessExpTime  = time.Now().Add(24 * time.Hour)
		refreshExpTime = time.Now().AddDate(0, 1, 0)
	)

	accessClaims := &m.JWTClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(t.key))
	if err != nil {
		return m.JWT{}, fmt.Errorf("GenerateToken (1): %w", err)
	}

	refreshClaims := &m.JWTClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshToken, err := refreshToken.SignedString([]byte(t.key))
	if err != nil {
		return m.JWT{}, fmt.Errorf("GenerateToken (2): %w", err)
	}

	tokens := m.JWT{
		ID:           id,
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokens, nil
}

func (t *JWT) verifyToken(signedToken string) error {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.key), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return fmt.Errorf("VerifyToken (1): %w", ErrTokenExpired)
			}
		}
		return fmt.Errorf("VerifyToken (2): %w", err)
	}

	if !token.Valid {
		return errors.New("VerifyToken (3): invalid token")
	}

	return nil
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("getTokenFromHeader (1): authorization token is missing")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("getTokenFromHeader (2): invalid token format")
	}

	return tokenParts[1], nil
}
