package logger

import (
	"fmt"
	"os"
	"path/filepath"

	models "github.com/Wires-Solucao-e-Servicos/golang-logger-module/models"

	"github.com/pelletier/go-toml"
)

var ClientName string
var SMTPConfig *models.SMTP

func SetClientName(name string) {
	if name != "" {
		ClientName = name
		return
	}

	if envName := os.Getenv("CLIENT_NAME"); envName != "" {
		ClientName = envName
		return
	}

	ClientName = "Undefined"
}

func LoadSMTPConfig(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config *models.SMTP

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SMTP configuration: %w", err)
	}

	if  config.Server == "" ||
			config.Port == 0 ||
			config.Username == "" ||
			config.Password == "" ||
			config.From == "" ||
			len(config.To) == 0 {
		return fmt.Errorf("missing required SMTP configuration fields")
	}

	SMTPConfig = config
	
	return nil
}