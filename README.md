# Golang Logger Module

A comprehensive Go logging module with file-based logging and email notification capabilities, designed for the Watchdog Service monitoring system.

## Features

- **Thread-safe logging** with goroutine-based message processing
- **Multiple log levels**: Info, Warning, Debug, Error
- **Automatic email notifications** for errors via SMTP
- **File-based logging** with automatic directory creation
- **Caller information tracking** (file and line number)
- **Graceful shutdown** with proper resource cleanup
- **Configurable client naming** via environment variables or direct setting

## Installation

```bash
go get github.com/Wires-Solucao-e-Servicos/golang-logger-module
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/Wires-Solucao-e-Servicos/golang-logger-module/logger"
)

func main() {
    //Optional: Set client name ("Undefined" if not set here or by the "CLIENT_NAME" environment variable)
    logger.SetClientName("MyApp")

    //Optional: Set SMTP configuration (No notification service if not loaded)
    logger.LoadSMTPConfig("config/config.toml")

    // Log different levels
    logger.Info("APP_START", "MAIN", "Application started successfully")
    logger.Warning("CONFIG_DEFAULT", "MAIN", "Using default configuration")
    logger.Debug("USER_ACTION", "AUTH", "User login attempt")
    logger.Error("DB_CONNECTION", "DATABASE", fmt.Errorf("failed to connect to database"))
    
    // Graceful shutdown
    defer logger.Close()
}
```

## API Reference

### Core Functions

```go
logger.Info(code, module, text string)
```
Logs informational messages.

```go
logger.Warning(code, module, text string)
```
Logs warning messages.

```go
logger.Debug(code, module, text string)`
```
Logs debug messages.

```go
logger.Error(code, module string, err error)`
```
Logs error messages and automatically sends email notification if SMTP is configured.

```go
logger.Close()`
```
Gracefully shuts down the logger, ensuring all messages are written before termination.

### Configuration Functions

```go
logger.SetClientName(name string)
```
Sets the client name for log identification. Falls back to `CLIENT_NAME` environment variable if name is empty.

```go
logger.GetClientName() string
```
Returns the current client name (thread-safe).

```go
logger.LoadSMTPConfig(path string) error
```
Loads SMTP configuration from a TOML file for error notifications.

```go
logger.ValidateSMTPConfig(s *models.SMTP) error
```
Validate provided SMTP Configuration

```go
logger.GetSMTPConfig() *models.SMTP
```
Returns the current SMTP configuration (thread-safe).

### Adicional Functions

```go
logger.SendEmail(values models.Notification) error
```

Send email according to loaded configuration.

## Log Format

Logs follow this structured format:
```
[LEVEL] [TIMESTAMP] [MODULE] [CODE] [FILE:LINE] > message.
```

Example:
```
[ERR] [23/06/2025 14:30:15] [DATABASE] [CONN_FAIL] [main.go:42] > failed to connect to database.
```

## Directory Structure

The logger automatically creates the following directory structure:

**Windows:**
```
C:\Wires Workspace\Watchdog Service\Logs\Logs.txt
```

**Unix/Linux/macOS:**
```
~/Watchdog Service/Logs/Logs.txt
```

## SMTP Configuration Schema

```go
type SMTP struct {
    Server   string   `toml:"server"`   // SMTP Server Address
    Port     int      `toml:"port"`     // SMTP Server Port
    Username string   `toml:"username"` // SMTP Username
    Password string   `toml:"password"` // SMTP Password
    From     string   `toml:"from"`     // Sender Email Address
    To       []string `toml:"to"`       // Recipient Email Addresses
}
```

Load the SMTP configuration:

```go
func main() {
    // Load SMTP config for error notifications
    if err := logger.LoadSMTPConfig("config.toml"); err != nil {
        log.Printf("Failed to load SMTP config: %v", err)
    }
    
    // Your application code here
    logger.Error("DB_CONNECTION", "DATABASE", fmt.Errorf("failed to connect to database"))
}
```

## Error Handling

The module includes comprehensive error handling:

- **File operations**: Automatic directory creation with proper permissions
- **SMTP validation**: Email address validation and connection testing
- **Goroutine safety**: Protected channel operations and graceful shutdown
- **Resource cleanup**: Automatic file closing and goroutine termination

## Thread Safety

All public functions are thread-safe:
- Configuration access protected by `sync.RWMutex`
- Singleton logger instance with `sync.Once`
- Channel-based message processing prevents race conditions
- Graceful shutdown with `sync.WaitGroup`

## Environment Variables

- `CLIENT_NAME`: Used if not set through `SetClientName()`

## Dependencies

- `github.com/pelletier/go-toml` - TOML configuration parsing
- `github.com/jordan-wright/email` - Email sending functionality

## Examples

### Custom Error Handling

```go
func processData() error {
    if err := validateInput(); err != nil {
        logger.Error("INPUT_VALIDATION", "PROCESSOR", err)
        return fmt.Errorf("validation failed: %w", err)
    }
    
    logger.Info("PROCESS_START", "PROCESSOR", "Data processing initiated")
    
    if err := performProcessing(); err != nil {
        logger.Error("PROCESS_EXECUTION", "PROCESSOR", err)
        os.Exit(1)
    }
    
    logger.Info("PROCESS_SUCCESS", "PROCESSOR", "Data processed successfully")
    return nil
}
```

## License

This project is proprietary software owned by Wires Solução e Serviços.

## Author

Wires Solução e Serviços  
Email: vinicius@wires.com.br

---

© 2025 Wires Solução e Serviços. All rights reserved.
