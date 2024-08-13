package model

import "github.com/dgrijalva/jwt-go"

type Player struct {
	ID     string `json:"steamid"`
	Name   string `json:"personaname"`
	URL    string `json:"profileurl"`
	Avatar string `json:"avatarfull"`
}

type PlayerResponse struct {
	Response struct {
		Players []Player `json:"players"`
	} `json:"response"`
}

type Profile struct {
	ID           string `json:"id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	URL          string `json:"url" validate:"required"`
	Avatar       string `json:"avatar" validate:"required"`
	Kills        int    `json:"kills" validate:"required"`
	Deaths       int    `json:"deaths" validate:"required"`
	HeadshotRate int    `json:"headshot_rate" validate:"required"`
}

type Stats struct {
	Kills     int
	Deaths    int
	Headshots int
}

type JWTClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type JWT struct {
	ID           string `json:"id" validate:"required"`
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
