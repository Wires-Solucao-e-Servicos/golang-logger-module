<<<<<<< HEAD
# Golang Logger
=======
# Golang Logger Module
>>>>>>> 3a7b554 (docs: update README.md)

A comprehensive Go logging module with file-based logging and email notification capabilities, designed for the Watchdog Service monitoring system.

## Features

<<<<<<< HEAD
- **Multi-Level Logging**: Support for Info, Warning, Debug, and Error log levels
- **File-Based Persistence**: Automatic log file creation and management
- **Email Notifications**: Automatic email alerts for error-level events via SMTP
- **Concurrent Logging**: Thread-safe logging with goroutine-based message processing
- **Caller Information**: Automatic capture of file name and line number for log entries
- **Cross-Platform**: Works on Windows and Unix-based systems
- **Singleton Pattern**: Global logger instance with proper initialization
- **TOML Configuration**: Easy SMTP configuration management using TOML files

## Project Structure

```
golang-logger-module/
├── config/
│   ├── loader.go          # Configuration loading logic
│   └── config.toml        # SMTP configuration file
├── logger/
│   └── logger.go          # Main logging functionality
├── models/
│   └── types.go           # Data structures and types
├── notification/
│   └── notify.go          # Email notification functionality
├── go.mod                 # Go module definition
└── go.sum                 # Dependency checksums
```
=======
- **Thread-safe logging** with goroutine-based message processing
- **Multiple log levels**: Info, Warning, Debug, Error
- **Automatic email notifications** for errors via SMTP
- **File-based logging** with automatic directory creation
- **Caller information tracking** (file and line number)
- **Graceful shutdown** with proper resource cleanup
- **Configurable client naming** via environment variables or direct setting
>>>>>>> 3a7b554 (docs: update README.md)

## Installation

```bash
<<<<<<< HEAD
go get github.com/Wires-Solucao-e-Servicos/golang-logger-module/log
```

## Configuration

### SMTP Configuration

Add the `config/config.toml` file with your SMTP server details for error notifications:

```toml
[smtp]
server = "smtp.gmail.com"
port = 587
username = "your-email@gmail.com"
password = "your-app-password"
from = "your-email@gmail.com"
to = ["recipient1@example.com", "recipient2@example.com"]
```

### Log File Location

The logger automatically creates log files in the following locations:

- **Windows**: `C:\Wires Workspace\Watchdog Service\Logs\Logs.txt`
- **Unix/Linux**: `~/Watchdog Service/Logs/Logs.txt`

## Usage

### Basic Logging
=======
go get github.com/Wires-Solucao-e-Servicos/golang-logger-module
```

## Quick Start

### Basic Usage
>>>>>>> 3a7b554 (docs: update README.md)

```go
package main

import (
    "fmt"
<<<<<<< HEAD
    
    "github.com/Wires-Solucao-e-Servicos/golang-logger-module/log"
)

func main() {
    // Initialize logger (called automatically on first use)
    
    // Info level logging
    logger.Info("APP_START", "MAIN", "Application started successfully")
    
    // Warning level logging
    logger.Warning("LOW_DISK", "SYSTEM", "Disk space is running low")
    
    // Debug level logging
    logger.Debug("DB_QUERY", "DATABASE", "Executing user query")
    
    // Error level logging (automatically sends email notification)
    err := fmt.Errorf("database connection failed")
    logger.Error("DB_ERROR", "DATABASE", err)
    
    // Properly close logger before application exit
=======
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
>>>>>>> 3a7b554 (docs: update README.md)
    defer logger.Close()
}
```

<<<<<<< HEAD
### Log Format

All log entries follow this standardized format:

=======
## API Reference

### Core Functions

#### `logger.Info(code, module, text string)`
Logs informational messages.

#### `logger.Warning(code, module, text string)`
Logs warning messages.

#### `logger.Debug(code, module, text string)`
Logs debug messages.

#### `logger.Error(code, module string, err error)`
Logs error messages and automatically sends email notification if SMTP is configured.

#### `logger.Close()`
Gracefully shuts down the logger, ensuring all messages are written before termination.

### Configuration Functions

#### `logger.SetClientName(name string)`
Sets the client name for log identification. Falls back to `CLIENT_NAME` environment variable if name is empty.

#### `logger.GetClientName() string`
Returns the current client name (thread-safe).

#### `logger.LoadSMTPConfig(path string) error`
Loads SMTP configuration from a TOML file for error notifications.

### `logger.ValidateSMTPConfig(s *models.SMTP) error`
Validate provided SMTP Configuration

#### `logger.GetSMTPConfig() *models.SMTP`
Returns the current SMTP configuration (thread-safe).

### Adicional Functions

### `logger.SendEmail(values models.Notification) error`
Send email according to loaded configuration.

## Log Format

Logs follow this structured format:
>>>>>>> 3a7b554 (docs: update README.md)
```
[LEVEL] [TIMESTAMP] [MODULE] [CODE] [FILE:LINE] > message.
```

<<<<<<< HEAD
Example log entries:
```
[INF] [15/01/2024 10:30:15] [MAIN] [APP_START] [main.go:10] > application started successfully.
[WRG] [15/01/2024 10:31:20] [SYSTEM] [LOW_DISK] [monitor.go:45] > disk space is running low.
[ERR] [15/01/2024 10:32:10] [DATABASE] [CONN_FAIL] [db.go:23] > database connection failed.
```

### Logging Functions

| Function | Level | Email Alert | Description |
|----------|-------|-------------|-------------|
| `Info(code, module, text)` | INFO | No | General information messages |
| `Warning(code, module, text)` | WARNING | No | Warning messages for potential issues |
| `Debug(code, module, text)` | DEBUG | No | Debug information for development |
| `Error(code, module, err)` | ERROR | Yes | Error messages with automatic email notification |

### Advanced Usage

```go
package main

import (
    "database/sql"
    "fmt"
    
    "github.com/Wires-Solucao-e-Servicos/golang-logger-module/log"
)

func connectDatabase() error {
    logger.Info("DB_INIT", "DATABASE", "Initializing database connection")
    
    // Simulate database connection
    db, err := sql.Open("mysql", "connection_string")
    if err != nil {
        logger.Error("DB_CONN_FAIL", "DATABASE", err)
        return err
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        logger.Error("DB_PING_FAIL", "DATABASE", err)
        return err
    }
    
    logger.Info("DB_SUCCESS", "DATABASE", "Database connected successfully")
    return nil
}

func main() {
    // Application startup
    logger.Info("APP_INIT", "MAIN", "Starting Watchdog Service")
    
    // Connect to database
    if err := connectDatabase(); err != nil {
        logger.Error("STARTUP_FAIL", "MAIN", fmt.Errorf("failed to start application: %w", err))
        return
    }
    
    logger.Info("APP_READY", "MAIN", "Watchdog Service is ready")
    
    // Ensure proper cleanup
    defer func() {
        logger.Info("APP_SHUTDOWN", "MAIN", "Shutting down Watchdog Service")
        logger.Close()
    }()
}
```

## Key Features

### Automatic Initialization

The logger uses a singleton pattern with automatic initialization. You don't need to manually initialize it - it's ready to use on the first logging call.

### Concurrent Safety

The logger uses goroutines and channels for thread-safe logging operations, ensuring no log messages are lost in concurrent environments.

### Error Notifications

Error-level logs automatically trigger email notifications with the following information:
- Timestamp of the error
- Error code
- Location (file and line where error occurred)
- Full formatted log message

### Caller Information

The logger automatically captures and includes the filename and line number where the log function was called, making debugging easier.

### Cross-Platform Support

Automatically detects the operating system and creates log files in appropriate locations for both Windows and Unix-based systems.

## Dependencies

- [github.com/pelletier/go-toml](https://github.com/pelletier/go-toml) - TOML configuration parsing
- [github.com/jordan-wright/email](https://github.com/jordan-wright/email) - Email sending functionality

## Best Practices

### Error Handling
```go
// Good: Use specific error codes and modules
logger.Error("DB_TIMEOUT", "DATABASE", err)
logger.Error("API_RATE_LIMIT", "HTTP_CLIENT", err)

// Bad: Generic error codes
logger.Error("ERROR", "SYSTEM", err)
```

### Module Organization
```go
// Use consistent module names throughout your application
logger.Info("INIT", "AUTH_SERVICE", "Authentication service starting")
logger.Info("USER_LOGIN", "AUTH_SERVICE", "User logged in successfully")
logger.Error("AUTH_FAIL", "AUTH_SERVICE", err)
```

### Proper Cleanup
```go
func main() {
    // Your application logic here
    
    // Always close logger before exit
    defer logger.Close()
}
```

## Environment Configuration

The module includes a `CLIENT_NAME` variable that defaults to "Development". Modify this in the config package for scenarios:

```go
var CLIENT_NAME string = "Customer 1"
```

## Troubleshooting

### Common Issues

1. **Permission Errors**: Ensure the application has write permissions to the log directory
2. **SMTP Failures**: Check your SMTP configuration and network connectivity
3. **Log File Not Created**: Verify disk space and directory permissions

### Log File Permissions

The logger creates files with `0644` (read and write the file or directory and other users can only read it) permissions and directories with `0755` (read, write, and execute permissions) permissions. Ensure your application runs with appropriate user permissions.

## Security Considerations

- Store SMTP credentials securely (consider environment variables for production)
- Use app-specific passwords for Gmail and similar providers
- Ensure log files don't contain sensitive information
- Consider log rotation for long-running applications
- Use TLS/SSL for SMTP connections
=======
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
>>>>>>> 3a7b554 (docs: update README.md)

## Author

Wires Solução e Serviços  
Email: vinicius@wires.com.br

---

© 2025 Wires Solução e Serviços. All rights reserved.
