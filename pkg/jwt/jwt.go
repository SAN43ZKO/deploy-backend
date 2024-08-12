package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	m "github.com/cs2-server/backend/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if err := t.verifyToken(tokenString); err != nil {
			logrus.Errorln("JWT (2): ", err)
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (t *JWT) GenerateTokens(id string) (m.JWT, error) {
	var (
		accessExpTime  = time.Now().Add(15 * time.Minute)
		refreshExpTime = time.Now().Add(24 * time.Hour)
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
	token, err := jwt.ParseWithClaims(signedToken, &m.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("VerifyToken (1): unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.key), nil
	})
	if err != nil {
		return fmt.Errorf("VerifyToken (2): %w", err)
	}

	if !token.Valid {
		return errors.New("VerifyToken (3): invalid token")
	}

	claims, ok := token.Claims.(*m.JWTClaims)
	if !ok {
		return errors.New("VerifyToken (4): invalid token claims")
	}

	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		return errors.New("VerifyToken (5): token has expired")
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
