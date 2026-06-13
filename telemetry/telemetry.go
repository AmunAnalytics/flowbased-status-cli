package telemetry

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"main/config"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IsDisabled() bool {
	return config.GetConfig().General.DisableTelemetry
}

var TelemetryMagicString = "custom"

func GetUid() (string, error) {
	if IsDisabled() {
		return "", fmt.Errorf("Telemetry is disabled")
	}

	pathFolder := config.GetConfigFolder()
	path := filepath.Join(pathFolder, "amun-cli-telemetry")
	os.MkdirAll(pathFolder, os.ModePerm)
	data, err := os.ReadFile(path)
	if err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	secret := make([]byte, 16)
	for i := range secret {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		secret[i] = charset[n.Int64()]
	}

	s := string(secret)
	if err := os.WriteFile(path, []byte(s), 0600); err != nil {
		return "", err
	}
	return s, nil
}

func Register(cli_type string, version string) {
	if IsDisabled() {
		return
	}

	uid, _ := GetUid()
	payload, err := json.Marshal(map[string]string{
		"uid":     uid,
		"version": version,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	})
	if err != nil {
		return
	}

	var url string
	if version == "dev" {
		url = "http://localhost:8080/telemetry/visit/%s/"
	} else {
		url = "https://ui.dl.amunanalytics.eu/telemetry/visit/%s/"
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(url, cli_type),
		bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Add("X-TELEMETRY-SECRET", TelemetryMagicString)

	client.Do(req)
}
