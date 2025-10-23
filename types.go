package logger

type SMTP struct {
	Server   string   `toml:"server" env:"SMTP_HOST"`
	Port     int      `toml:"port" env:"SMTP_PORT"`
	Username string   `toml:"username" env:"SMTP_USERNAME"`
	Password string   `toml:"password" env:"SMTP_PASSWORD"`
	From     string   `toml:"from" env:"SMTP_FROM"`
	To       []string `toml:"to" env:"SMTP_TO"`
}

type Notification struct {
	Datetime string
	Code     string
	Location string
	Details  string
}