package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var oauthConfig *oauth2.Config
var oauthStateStore = make(map[string]bool)

type TwitterUser struct {
	Data struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

func InitOAuth() {
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("TWITTER_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BACKEND_URL") + "/api/auth/twitter/callback",
		Scopes:       []string{"tweet.read", "users.read", "bookmark.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
			TokenURL: "https://api.twitter.com/2/oauth2/token",
		},
	}
}

func GenerateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	oauthStateStore[state] = true
	return state
}

func ValidateStateToken(state string) bool {
	valid, exists := oauthStateStore[state]
	if exists {
		delete(oauthStateStore, state)
	}
	return valid
}

func GetAuthURL(state string) string {
	return oauthConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", "challenge"),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)
}

func ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := oauthConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", "challenge"),
	)
	return token, err
}

func GetTwitterUserInfo(ctx context.Context, token *oauth2.Token) (*TwitterUser, error) {
	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.twitter.com/2/users/me?user.fields=profile_image_url")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var twitterUser TwitterUser
	err = json.NewDecoder(resp.Body).Decode(&twitterUser)
	if err != nil {
		return nil, err
	}

	return &twitterUser, nil
}
