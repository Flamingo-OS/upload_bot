package core

type OTA struct {
	Version             string            `json:"version"`
	Date                string            `json:"date"`
	DownloadSources     map[string]string `json:"download_sources"`
	FileName            string            `json:"file_name"`
	FileSize            string            `json:"file_size"`
	ShaSum              string            `json:"sha_512"`
	PreBuildIncremental string            `json:"pre_build_incremental"`
}
