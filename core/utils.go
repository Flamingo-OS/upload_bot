package core

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFile(zipFilePath string, fileName string) (string, error) {
	dst := DumpPath + fileName
	archive, err := zip.OpenReader(zipFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		Log.Info("Extracting %s", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			Log.Error("illegal file path: %s", filePath)
			return "", fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			Log.Error("Error while creating directory: %s", err)
			return "", err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			Log.Error("Error while opening file: %s", err)
			return "", err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			Log.Error("Error while opening file in archive: %s", err)
			return "", err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			Log.Error("Error while copying file: %s", err)
			return "", err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return "", nil
}
