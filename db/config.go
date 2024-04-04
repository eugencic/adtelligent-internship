package db

type Config struct {
	host     string
	port     string
	user     string
	password string
	database string
}

var dbConfig = Config{
	host:     "localhost",
	port:     "3306",
	user:     "user",
	password: "password",
	database: "database",
}
