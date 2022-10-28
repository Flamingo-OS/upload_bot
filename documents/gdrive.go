package documents

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Flamingo-OS/upload-bot/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
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

	core.DriveClient = client
	core.DriveService = srv
	return nil
}

func getFileFromId(d *drive.Service, fileId string) (*drive.File, error) {
	core.Log.Info("Finding details about the file to download")
	f, err := d.Files.Get(fileId).SupportsAllDrives(true).SupportsTeamDrives(true).Do()
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return nil, err
	}
	core.Log.Info("Title: %v \n", f.Title)
	return f, nil
}

// DownloadFile downloads the content of a given file object
func downloadFile(d *drive.Service, t http.RoundTripper, f *drive.File) (string, error) {
	core.Log.Info("Initialising download")
	downloadUrl := f.DownloadUrl
	title := f.Title
	if downloadUrl == "" {
		core.Log.Errorln("An error occurred: File is not downloadable")
		return "", nil
	}
	core.Log.Info("Fetching requests")
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		core.Log.Errorf("error occurred: %v\n", err)
		return "", err
	}
	core.Log.Info("Fetching response body")
	resp, err := t.RoundTrip(req)
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	return downloadSaver(resp, title)
}

func parseFileId(url string) string {
	re, _ := regexp.Compile(`[-\w]{25,}`)
	return re.FindString(url)
}

func gDriveDownloader(url string) (string, error) {
	fileId := parseFileId(url)
	f, err := getFileFromId(core.DriveService, fileId)
	if err != nil {
		core.Log.Errorf("An error occurred: %v\n", err)
		return "", err
	}
	return downloadFile(core.DriveService, core.DriveClient.Transport, f)
}
