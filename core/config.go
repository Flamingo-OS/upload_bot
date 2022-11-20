package core

import (
	"encoding/json"
	"os"
)

type BotConfig struct {
	BotToken                string `json:"bot_token"`
	MongoDbConnectionString string `json:"connection_string"`
	GDriveClientId          string `json:"gdrive_client_id"`
	GDriveClientSecret      string `json:"gdrive_client_secret"`
	GDriveAccessToken       string `json:"gdrive_access_token"`
	GDriveRefreshToken      string `json:"gdrive_refresh_token"`
	GDriveExpiry            string `json:"gdrive_expiry"`
	OneDriveClientId        string `json:"onedrive_client_id"`
	OneDriveClientSecret    string `json:"onedrive_client_secret"`
	OneDriveTenantId        string `json:"onedrive_tenant_id"`
	OneDriveRefreshToken    string `json:"onedrive_refresh_token"`
	GithubToken             string `json:"github_token"`
}

func NewBotConfig(fileName string) *BotConfig {
	ac := &BotConfig{}
	b, err := os.ReadFile(fileName)
	if err != nil {
		Log.Error(err)
		Log.Info("Reading config for docker container")
		b, err := os.ReadFile("/etc/secrets/config.json") // for docker
		if err != nil {
			Log.Error(err)
			Log.Info("Reading config from environment")
			ac = configFromEnv()
			return ac
		}
		err = json.Unmarshal(b, ac)
		if err != nil {
			Log.Fatal(err)
		}
		return ac
	}
	err = json.Unmarshal(b, &ac)
	if err != nil {
		Log.Fatal(err)
	}
	return ac
}

func configFromEnv() *BotConfig {
	ac := &BotConfig{}
	var isPresent bool
	ac.BotToken, isPresent = os.LookupEnv("BOT_TOKEN")
	checkIfPresent(isPresent, "BOT_TOKEN")
	ac.MongoDbConnectionString, isPresent = os.LookupEnv("MONGO_DB_CONNECTION_STRING")
	checkIfPresent(isPresent, "MONGO_DB_CONNECTION_STRING")
	ac.GDriveClientId, isPresent = os.LookupEnv("GDRIVE_CLIENT_ID")
	checkIfPresent(isPresent, "GDRIVE_CLIENT_ID")
	ac.GDriveClientSecret, isPresent = os.LookupEnv("GDRIVE_CLIENT_SECRET")
	checkIfPresent(isPresent, "GDRIVE_CLIENT_SECRET")
	ac.GDriveExpiry, isPresent = os.LookupEnv("GDRIVE_EXPIRY")
	checkIfPresent(isPresent, "GDRIVE_EXPIRY")
	ac.GDriveAccessToken, isPresent = os.LookupEnv("GDRIVE_ACCESS_TOKEN")
	checkIfPresent(isPresent, "GDRIVE_ACCESS_TOKEN")
	ac.GDriveRefreshToken, isPresent = os.LookupEnv("GDRIVE_REFRESH_TOKEN")
	checkIfPresent(isPresent, "GDRIVE_REFRESH_TOKEN")
	ac.OneDriveClientId, isPresent = os.LookupEnv("ONEDRIVE_CLIENT_ID")
	checkIfPresent(isPresent, "ONEDRIVE_CLIENT_ID")
	ac.OneDriveClientSecret, isPresent = os.LookupEnv("ONEDRIVE_CLIENT_SECRET")
	checkIfPresent(isPresent, "ONEDRIVE_CLIENT_SECRET")
	ac.OneDriveTenantId, isPresent = os.LookupEnv("ONEDRIVE_TENANT_ID")
	checkIfPresent(isPresent, "ONEDRIVE_TENANT_ID")
	ac.OneDriveRefreshToken, isPresent = os.LookupEnv("ONEDRIVE_REFRESH_TOKEN")
	checkIfPresent(isPresent, "ONEDRIVE_REFRESH_TOKEN")
	ac.GithubToken, isPresent = os.LookupEnv("GITHUB_TOKEN")
	checkIfPresent(isPresent, "GITHUB_TOKEN")
	return ac
}

func checkIfPresent(isPresent bool, envVar string) {
	if !isPresent {
		Log.Fatalf("%s is not present in the environment", envVar)
	}
}
