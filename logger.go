package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Log struct {
	wg		  sync.WaitGroup
	quit    chan struct{}
	file    *os.File
	logger  *log.Logger
	channel chan string
}

var (
	instance   			 *Log
	once       			 sync.Once
	defaultDirectory string
)

func GetCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}

	return filepath.Base(file), line
}

func FormatLog(level, module, code, text string) string {

	file, line := GetCallerInfo()

	return fmt.Sprintf("[%s] [%s] [%s] [%s] [%s:%d] > %s.", level, Timestamp(), module, code, file, line, strings.ToLower(text))
}

func Timestamp() string {
	return time.Now().Format("02/01/2006 15:04:05")
}

func SetLoggerDirectory(path string) error {
	rwmu.Lock()
	defer rwmu.Unlock()

	if strings.TrimSpace(path) == "" {
			return fmt.Errorf("empty path: logger directory not changed")
	}

	defaultDirectory = path
	return nil
}

func CreateLoggerDirectory() (*os.File, error) {
	folderName := "Logger Module"
	var baseDirectory string

	if defaultDirectory != "" {
			baseDirectory = defaultDirectory
	} else {
			if runtime.GOOS == "windows" {
					baseDirectory = "C:\\Project"
			} else {
					home, err := os.UserHomeDir()
					if err != nil {
							return nil, fmt.Errorf("home dir not found: %w", err)
					}
					baseDirectory = home
			}
	}

	programDirectory := filepath.Join(baseDirectory, folderName)

	if err := os.MkdirAll(programDirectory, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	logPath := filepath.Join(programDirectory, "Logs.txt")

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
			header := fmt.Sprintf("[Logger Module - %s]\n\n", GetClientName())
			if err := os.WriteFile(logPath, []byte(header), 0644); err != nil {
					return nil, fmt.Errorf("failed to create log file: %w", err)
			}
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return logFile, nil
}

func Init() {

	once.Do(func() {

		logFile, err := CreateLoggerDirectory()
		if err != nil {
			log.Fatal("failed to create logger directory: %w", err)
		}

		instance = &Log{
			quit: 	 make(chan struct{}),
			file:    logFile,
			logger:  log.New(logFile, "", 0),
			channel: make(chan string, 100),
		}

		instance.wg.Add(1)
		go instance.run()

		instance.logger.Printf("\n\n[%s] [%s] [%s] [%s] [%s:%d] > %s.\n", "LOG", Timestamp(), "LOGGER", "LOGGER_INIT", "logger.go", 0, "logger initialized")
	})
}

func (l *Log) run() {

	defer l.wg.Done()
	defer l.file.Close()

	for {
		select {
			case message := <-l.channel:
				l.logger.Println(message)
			case <- l.quit:
				close(l.channel)
				for message := range l.channel {
					l.logger.Println(message)
				}
				return
		}
	}
}

func Close() {
	if instance != nil {
		instance.quit <- struct{}{}
		instance.wg.Wait()
	}
}

func Info(code, module, text string) {
	Init()

	message := FormatLog("INF", module, code, text)
	instance.channel <- message
}

func Warning(code, module, text string) {
	Init()

	message := FormatLog("WRG", module, code, text)
	instance.channel <- message
}

func Debug(code, module, text string) {
	Init()

	message := FormatLog("DBG", module, code, text)
	instance.channel <- message
}

func Error(code, module string, err error) {
	Init()
	
	text := err.Error()

	message := FormatLog("ERR", module, code, text)
	instance.channel <- message

	instance.wg.Add(1)
	go func() {

		defer instance.wg.Done()

		if instance == nil || SMTPConfig == nil {
			return
		}
		
		notification := Notification{
			Datetime: Timestamp(),
			Code:     code,
			Location: fmt.Sprint(GetCallerInfo()),
			Details:  message,
		}
		
		var errMsg string

		if err := SendEmail(notification); err != nil {
			errMsg = FormatLog("ERR", "LOGGER_NOTIFY", "SMTP_ERROR", fmt.Sprintf("failed to send notification: %v", message))
		} 

		if errMsg != "" {
			instance.channel <- errMsg
		}

	} ()
	
}
