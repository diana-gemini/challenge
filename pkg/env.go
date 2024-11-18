package pkg

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	PORT          string `mapstructure:"PORT"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      int    `mapstructure:"SMTP_PORT"`
	SMTPUser      string `mapstructure:"SMTP_USER"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env: ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
