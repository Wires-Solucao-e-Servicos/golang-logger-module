package logger

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"sync"

	models "github.com/Wires-Solucao-e-Servicos/golang-logger-module/models"

	"github.com/pelletier/go-toml"
)

var (
	rwmu 				sync.RWMutex
	SMTPConfig 	*models.SMTP
)

var clientName string = "Undefined"

func SetClientName(name string) {
	rwmu.Lock()
	defer rwmu.Unlock()

	if name != "" {
		clientName = name
		return
	}

	if envName := os.Getenv("CLIENT_NAME"); envName != "" {
		clientName = envName
		return
	}
}

func GetClientName() string {
	return clientName
}

func ValidateSMTPConfig(s *models.SMTP) error {
	if s.Server == "" {
		return fmt.Errorf("server is required")
	}
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if s.Username == "" {
		return fmt.Errorf("username is required")
	}
	if s.Password == "" {
		return fmt.Errorf("password is required")
	}
	if s.From == "" {
		return fmt.Errorf("from address is required")
	}
	if len(s.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	
	if _, err := mail.ParseAddress(s.From); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	
	for _, addr := range s.To {
		if _, err := mail.ParseAddress(addr); err != nil {
			return fmt.Errorf("invalid to address %s: %w", addr, err)
		}
	}
	
	return nil
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

	err = ValidateSMTPConfig(config)
	if err != nil {
		return fmt.Errorf("invalid SMTP config: %w", err)
	}

	rwmu.Lock()
	SMTPConfig = config

	rwmu.Unlock()
	
	return nil
}

func GetSMTPConfig() *models.SMTP {
	rwmu.Lock()
	defer rwmu.Unlock()

	return SMTPConfig
}