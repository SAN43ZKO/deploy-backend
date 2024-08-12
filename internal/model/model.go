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
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Avatar       string `json:"avatar"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	HeadshotRate int    `json:"headshot_rate"`
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
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
