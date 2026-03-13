package env

import "os"

type Env string

const (
	ENV         Env = "ENV"
	APP_PORT    Env = "APP_PORT"
	DB_HOST     Env = "DB_HOST"
	DB_PORT     Env = "DB_PORT"
	DB_USER     Env = "DB_USER"
	DB_PASSWORD Env = "DB_PASSWORD"
	DB_NAME     Env = "DB_NAME"
)

func (e Env) GetValue() string {
	return os.Getenv(string(e))
}
