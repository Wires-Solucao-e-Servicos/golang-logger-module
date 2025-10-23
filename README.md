# Golang Logger Module

A comprehensive Go logging module with file-based logging and email notification capabilities.

## Features

- **Thread-safe logging** with goroutine-based message processing
- **Multiple log levels**: Info, Warning, Debug, Error
- **Automatic email notifications** for errors via SMTP
- **File-based logging** with automatic directory creation
- **Caller information tracking** (file and line number)
- **Graceful shutdown** with proper resource cleanup
- **Flexible configuration**: Load from TOML file or environment variables

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
    "github.com/Wires-Solucao-e-Servicos/golang-logger-module"
)

func main() {
    // Optional: Set client name for log identification 
    // Defaults to "Undefined" if not set here or via the "CLIENT_NAME" environment variable
    logger.SetClientName("MyApp")

    // Optional: Load SMTP configuration from TOML or ENV
    // If not loaded, errors will only be logged locally
    if err := logger.LoadTOMLConfig("config/config.toml"); err != nil {
        log.Printf("Failed to load TOML config: %v", err)
    }
    // OR load from .env file
    // if err := logger.LoadENVConfig(); err != nil {
    //     log.Printf("Failed to load ENV config: %v", err)
    // }

    //Optional: Set custom logger directory
    // Defaults to "C:/Project" on Windows and the user's home dir elsewhere
    if err := logger.SetLoggerDirectory(""); err != nil {
        log.Printf("Failed to change default logger directory: %v", err)
    }

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
logger.Info(code, module, text string) //Logs informational messages.
```

```go
logger.Warning(code, module, text string) // Logs warning messages.
```

```go
logger.Debug(code, module, text string) // Logs debug messages.
```

```go
logger.Error(code, module string, err error) //Logs error messages and automatically sends email notification.
```

```go
logger.Close() // Shuts down the logger, ensuring all messages are written before termination.
```

### Configuration Functions

```go
logger.SetClientName(name string) // Sets the client name for log identification
```

```go
logger.GetClientName() string // Returns the current client name (thread-safe)
```

```go
logger.LoadTOMLConfig(path string) error // Loads SMTP configuration from a TOML file
```

```go
logger.LoadENVConfig() error // Loads SMTP configuration from environment variables (.env file)
```

```go
logger.ValidateSMTPConfig(s *SMTP) error // Validates provided SMTP configuration
```

```go
logger.GetSMTPConfig() *SMTP // Returns the current SMTP configuration (thread-safe)
```

```go
logger.SetLoggerDirectory(path string) error // Changes the default logger directory
```

### Adicional Functions

```go
logger.SendEmail(values Notification) error // Send email according to loaded configuration.
```

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

## Models

The module defines two main data structures:

### SMTP Configuration

The `SMTP` struct defines the email server configuration:

```go
type SMTP struct {
    Server   string   `toml:"server" env:"SMTP_SERVER"`     // SMTP server address (e.g., "smtp.gmail.com")
    Port     int      `toml:"port" env:"SMTP_PORT"`         // SMTP server port (e.g., 587 for TLS)
    Username string   `toml:"username" env:"SMTP_USERNAME"` // SMTP authentication username
    Password string   `toml:"password" env:"SMTP_PASSWORD"` // SMTP authentication password
    From     string   `toml:"from" env:"SMTP_FROM"`         // Sender email address
    To       []string `toml:"to" env:"SMTP_TO"`             // List of recipient email addresses
}
```

This struct is automatically populated when loading configuration via `LoadTOMLConfig()` or `LoadENVConfig()` and includes tags for both TOML and environment variable parsing.

### Notification Structure

The `Notification` struct represents an email notification payload:

```go
type Notification struct {
    Datetime string  // Timestamp when the event occurred
    Code     string  // Event identification code
    Location string  // Source file and line number where the event was logged
    Details  string  // Complete formatted log message with full context
}
```

This struct is used internally by the `Error()` function and can be used directly with `SendEmail()` for custom notifications.

## SMTP Configuration

To enable email notifications on errors, you have two options:

### Option 1: TOML Configuration File

Create a `config.toml` file with the SMTP configuration:

```toml
server = "smtp.gmail.com"
port = 587
username = "email@gmail.com"
password = "your-app-password"
from = "email@gmail.com"
to = ["recipient1@gmail.com", "recipient2@gmail.com"]
```

Load the configuration:

```go
if err := logger.LoadTOMLConfig("config/config.toml"); err != nil {
    log.Printf("Failed to load SMTP config: %v", err)
}
```

### Option 2: Environment Variables

Create a `.env` file in your project root:

```env
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=email@gmail.com
SMTP_TO=recipient1@gmail.com,recipient2@gmail.com
```

Load the configuration:

```go
if err := logger.LoadENVConfig(); err != nil {
    log.Printf("Failed to load ENV config: %v", err)
}
```

**Note:** For `SMTP_TO`, use comma-separated email addresses.

If neither configuration is loaded, errors will be logged locally but no email alerts will be sent.

## Email Notifications

You can send email notifications without logging an error by directly using the `SendEmail` function with a properly formatted `Notification` struct.

Example of sending a custom notification email:

```go
func main() {
    notification := logger.Notification{
        Datetime: logger.Timestamp(),
        Code:     "CUSTOM_ALERT",
        Location: fmt.Sprint(logger.GetCallerInfo()),
        Details:  "Custom alert message details here",
    }

    if err := logger.SendEmail(notification); err != nil {
        errMsg := logger.FormatLog("ERR", "LOGGER_NOTIFY", "SMTP_ERROR", fmt.Sprintf("failed to send notification: %v", err))
    }
}
```

Make sure the `Notification` fields follow the expected format so the email content and log formatting work correctly.

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
- `SMTP_SERVER`: SMTP server address
- `SMTP_PORT`: SMTP server port (defaults to 587 if invalid)
- `SMTP_USERNAME`: SMTP authentication username
- `SMTP_PASSWORD`: SMTP authentication password
- `SMTP_FROM`: Sender email address
- `SMTP_TO`: Comma-separated list of recipient email addresses

## Dependencies

- `github.com/pelletier/go-toml` - TOML configuration parsing
- `github.com/jordan-wright/email` - Email sending functionality
- `github.com/joho/godotenv` - Environment variable loading from .env files

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
