package config

type Default struct {
	Server ServerConfig
	//BotToken  string `env:"TELEGRAM_BOT_TOKEN,required"`
	//JWTSecret string `env:"JWT_SECRET,required"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
}
