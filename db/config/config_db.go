package config

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var DbConfig = Config{
	Host:     "localhost",
	Port:     "3306",
	User:     "user",
	Password: "password",
	Database: "database",
}
