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

	models "github.com/Wires-Solucao-e-Servicos/golang-logger-module/models"
)

type Log struct {
	wg		  sync.WaitGroup
	quit    chan struct{}
	file    *os.File
	logger  *log.Logger
	channel chan string
}

var (
	instance   *Log
	once       sync.Once
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

func CreateLoggerDirectory() (*os.File, error) {

	var programDirectory string
		folderName := "Watchdog Service"

		if runtime.GOOS == "windows" {
			programDirectory = filepath.Join("C:\\Wires Workspace", folderName)
		} else {
			homeDirectory, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("home dir not found: %s", err)
			}
			programDirectory = filepath.Join(homeDirectory, folderName)
		}

		logDirectory := filepath.Join(programDirectory, "Logs")
		if err := os.MkdirAll(logDirectory, 0755); err != nil {
			return nil, fmt.Errorf("failed to create logs dir: %v", err)
		}

		logPath := filepath.Join(logDirectory, "Logs.txt")

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			err = os.WriteFile(logPath, fmt.Appendf(nil, "[Wires Watchdog Service - %s] \n\n", GetClientName()), 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to create log file: %v", err)
			}
		}

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
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

	go func() {

		if instance == nil {
			return
		}
		
		notification := models.Notification{
			Datetime: Timestamp(),
			Code:     code,
			Location: fmt.Sprint(GetCallerInfo()),
			Details:  message,
		}
		
		if err := SendEmail(notification); err != nil {
			errMsg := FormatLog("ERR", "LOGGER_NOTIFY", "SMTP_ERROR", fmt.Sprintf("failed to send notification: %v", message))
			instance.channel <- errMsg
		} 
	} ()
	
}
