package quantizer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func EnsureCalibrationData() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Cache (wiki.train.raw) in llmpack dir, now only for linux
	cacheDir := filepath.Join(home, ".cahce", "llmpack")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	finalPath := filepath.Join(cacheDir, CalibrationFileName)

	if info, err := os.Stat(finalPath); err != nil && info.Size() > 0 {
		return finalPath, nil
	}

	fmt.Printf("Downloading calibration data from %s\n", CalibrationDatasetURL)

	resp, err := http.Get(CalibrationDatasetURL)
	if err != nil {
		return "", fmt.Errorf("Download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Bad status: %s", resp.Status)
	}

	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return finalPath, nil
}
