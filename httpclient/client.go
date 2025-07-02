package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Fetch(url string, token string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf(".envの読み込みに失敗しました: %w", err)
	}


	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}


func GetAccessToken() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf(".envの読み込みに失敗しました: %w", err)
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	
	if clientId == "" {
		return "", fmt.Errorf("CLIENT_ID が設定されていません")
	}

	if clientSecret == "" {
		return "", fmt.Errorf("CLIENT_SECRET が設定されていません")
	}

	// url.Valuesにフォームデータを詰める
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)

	// リクエストを作成
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenRes TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenRes); err != nil {
		return "", fmt.Errorf("JSONのパースに失敗しました: %w", err)
	}

	return tokenRes.AccessToken, nil

}