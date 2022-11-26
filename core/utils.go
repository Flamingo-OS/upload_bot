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

func createFile(f *zip.File, filePath string, dst string) error {
	if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
		Log.Error("illegal file path:", filePath)
		return fmt.Errorf("invalid file path")
	}
	if f.FileInfo().IsDir() {
		Log.Infof("creating directory...")
		os.MkdirAll(filePath, os.ModePerm)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		Log.Error("Error while creating directory:", err)
		return err
	}

	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		Log.Error("Error while opening file:", err)
		return err
	}
	defer dstFile.Close()

	fileInArchive, err := f.Open()
	if err != nil {
		Log.Error("Error while opening file in archive:", err)
		return err
	}
	defer fileInArchive.Close()

	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		Log.Error("Error while copying file:", err)
		return err
	}

	return nil
}

// provide a filepath to extract a particular file only
func UnzipFile(zipFilePath string, fileName string, dumpPath string, fPath string) (string, error) {
	dst := dumpPath + fileName
	archive, err := zip.OpenReader(zipFilePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		if fPath != "" && filePath == (dumpPath+fileName+"/"+fPath) {
			Log.Info("Extracting to ", filePath)
			createFile(f, filePath, dst)
			break
		} else if fPath != "" {
			continue
		}
		Log.Info("Extracting to ", filePath)
		createFile(f, filePath, dst)
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
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f := strings.Split(fileName, "/")
		os.MkdirAll(strings.Join(f[:len(f)-1], "/"), os.ModePerm)
	}
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
