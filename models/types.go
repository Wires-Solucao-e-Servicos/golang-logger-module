package models

type SMTP struct {
	Server   string   `toml:"server"`
	Port     int      `toml:"port"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	From     string   `toml:"from"`
	To       []string `toml:"to"`
}

type Notification struct {
	Datetime string
	Code     string
	Location string
	Details  string
}