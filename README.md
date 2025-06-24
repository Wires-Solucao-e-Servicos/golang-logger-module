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
    "github.com/Wires-Solucao-e-Servicos/golang-logger-module"
)

func main() {
    // Optional: Set client name for log identification 
    // Defaults to "Undefined" if not set here or via the "CLIENT_NAME" environment variable
    logger.SetClientName("MyApp")

    //Optional: Set SMTP configuration 
    // If not loaded, errors will only be logged locally
    if err := logger.LoadSMTPConfig("config/config.toml"); err != nil {
        log.Printf("Failed to load SMTP config: %v", err)
    }

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
logger.GetClientName() string // Returns the current client name (thread-safe).
```

```go
logger.LoadSMTPConfig(path string) error // Loads SMTP configuration from a TOML file for error notifications.
```

```go
logger.ValidateSMTPConfig(s *SMTP) error // Validate provided SMTP Configuration.
```

```go
logger.GetSMTPConfig() *SMTP // Returns the current SMTP configuration (thread-safe).
```

```go
logger.SetLoggerDirectory(path string) error// Change the default logger directory
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
    Server   string   `toml:"server"`    // SMTP server address (e.g., "smtp.gmail.com")
    Port     int      `toml:"port"`      // SMTP server port (e.g., 587 for TLS)
    Username string   `toml:"username"`  // SMTP authentication username
    Password string   `toml:"password"`  // SMTP authentication password
    From     string   `toml:"from"`      // Sender email address
    To       []string `toml:"to"`        // List of recipient email addresses
}
```

This struct is automatically populated when loading configuration via `LoadSMTPConfig()` and includes TOML tags for proper file parsing.

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

## SMTP Configuration Schema

To enable email notifications on errors, you must create a config.toml file with the SMTP configuration and load it using the `LoadSMTPConfig(path string)` function.

Use config/config.toml.example as a reference for the required structure. If this configuration isn’t loaded, errors will be logged locally but no email alerts will be sent.

```toml
port = 587
password = "password"
server = "smtp.gmail.com"
from = "email@gmail.com"
to = ["email@gmail.com"]
username = "email@gmail.com"
```

Load the SMTP configuration:

```go
func main() {
    // Use the relative path to your config file, e.g., "config/config.toml"
    if err := logger.LoadSMTPConfig("config/config.toml"); err != nil {
        log.Printf("Failed to load SMTP config: %v", err)
    }

    // Your application code here
    logger.Error("DB_CONNECTION", "DATABASE", fmt.Errorf("failed to connect to database"))
}
```

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
