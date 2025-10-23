package logger

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml"
)

var (
	rwmu 				sync.RWMutex
	SMTPConfig 	*SMTP
	clientName string = "Golang Logger"
)

func SetClientName(name string) {
	rwmu.Lock()
	defer rwmu.Unlock()

	if name != "" {
		clientName = name
		return
	}
}

func GetClientName() string {
	envName := os.Getenv("CLIENT_NAME");
	
	rwmu.RLock()
	defer rwmu.RUnlock()

	if clientName != "Golang Logger" {
		return clientName
	}

	if envName != "" {
		return envName
	} 

	return clientName
}

func ValidateSMTPConfig(s *SMTP) error {
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

func LoadENVConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load env config: %w", err);
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587
	}

	var to []string
  toStr := os.Getenv("SMTP_TO")
  if toStr != "" {
    to = strings.Split(toStr, ",")
  }

	config := &SMTP{
		Server: os.Getenv("SMTP_HOST"),
		Port: port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From: os.Getenv("SMTP_FROM"),
		To: to,
	}

	err = ValidateSMTPConfig(config)
	if err != nil {
		return fmt.Errorf("invalid SMTP config: %w", err)
	}

	rwmu.Lock()
	defer rwmu.Unlock()

	SMTPConfig = config
	
	return nil
}

func LoadTOMLConfig(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &SMTP{}

	err = toml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SMTP configuration: %w", err)
	}

	err = ValidateSMTPConfig(config)
	if err != nil {
		return fmt.Errorf("invalid SMTP config: %w", err)
	}

	rwmu.Lock()
	defer rwmu.Unlock()

	SMTPConfig = config
	
	return nil
}

func GetSMTPConfig() *SMTP {
	rwmu.RLock()
	defer rwmu.RUnlock()

	return SMTPConfig
}
