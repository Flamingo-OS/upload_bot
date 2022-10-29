package core

import (
	"archive/zip"
	"crypto/sha512"
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
		Log.Info("Extracting", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			Log.Error("illegal file path:", filePath)
			return "", fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			Log.Error("Error while creating directory:", err)
			return "", err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			Log.Error("Error while opening file:", err)
			return "", err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			Log.Error("Error while opening file in archive:", err)
			return "", err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			Log.Error("Error while copying file:", err)
			return "", err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return dst, nil
}

func FindShaSum(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		Log.Error("Error while opening file:", err)
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		Log.Error("Error while copying file:", err)
		return "", fmt.Errorf("error while copying file")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func writeToFile(fileName string, content string) error {
	f, err := os.Create(fileName)
	if err != nil {
		Log.Error("Error while creating file:", err)
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		Log.Error("Error while writing to file:", err)
		return err
	}
	return nil
}
