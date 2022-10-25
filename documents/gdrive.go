package documents

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Flamingo-OS/upload-bot/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Expiry       string `json:"expiry"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokenJson := DriveToken{
		AccessToken:  core.Config.GDriveAccessToken,
		RefreshToken: core.Config.GDriveRefreshToken,
		TokenType:    "Bearer",
		Expiry:       core.Config.GDriveExpiry,
	}
	b, err := json.Marshal(tokenJson)
	if err != nil {
		log.Fatalf("Unable to marshal token: %v", err)
		return nil
	}
	tok := &oauth2.Token{}
	json.Unmarshal(b, &tok)
	return config.Client(context.Background(), tok)
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     core.Config.GDriveClientId,
		ClientSecret: core.Config.GDriveClientSecret,
		RedirectURL:  "http://localhost:8080",
		Scopes: []string{
			drive.DriveScope,
		},
		Endpoint: google.Endpoint,
	}
}

func NewGdrive() error {
	ctx := context.Background()
	client := getClient(getConfig())

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		core.Log.Fatalf("Unable to retrieve Drive client: %v", err)
		return err
	}

	core.DriveService = srv
	return nil
}
