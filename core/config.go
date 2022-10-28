package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type BotConfig struct {
	BotToken                string `json:"bot_token"`
	MongoDbConnectionString string `json:"connection_string"`
	GDriveClientId          string `json:"gdrive_client_id"`
	GDriveClientSecret      string `json:"gdrive_client_secret"`
	GDriveAccessToken       string `json:"gdrive_access_token"`
	GDriveRefreshToken      string `json:"gdrive_refresh_token"`
	GDriveTokenType         string `json:"gdrive_token_type"`
	GDriveExpiry            string `json:"gdrive_expiry"`
	OneDriveClientId        string `json:"onedrive_client_id"`
	OneDriveClientSecret    string `json:"onedrive_client_secret"`
	OneDriveTenantId        string `json:"onedrive_tenant_id"`
	OneDriveRefreshToken    string `json:"onedrive_refresh_token"`
}

func NewBotConfig(fileName string) *BotConfig {
	ac := &BotConfig{}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}
	return ac
}
