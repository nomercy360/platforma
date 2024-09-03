package config

type Default struct {
	Server        ServerConfig
	Bepaid        Bepaid
	ExternalURL   string `env:"EXTERNAL_URL,required"`
	Notifications Notifications
	//JWTSecret string `env:"JWT_SECRET,required"`
	DBPath string `env:"DB_PATH" envDefault:"./app.db"`
	WebURL string `env:"WEB_URL" envDefault:"http://localhost:3000"`
	PayPal PayPal
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
}

type PayPal struct {
	ClientID     string `env:"PAYPAL_CLIENT_ID,required"`
	ClientSecret string `env:"PAYPAL_CLIENT_SECRET,required"`
	LiveMode     bool   `env:"PAYPAL_LIVE_MODE" envDefault:"false"`
}

type Bepaid struct {
	ShopID    string `env:"BEPAID_SHOP_ID,required"`
	SecretKey string `env:"BEPAID_SECRET_KEY,required"`
	ApiURL    string `env:"BEPAID_API_URL" envDefault:"https://checkout.bepaid.by"`
	TestMode  bool   `env:"BEPAID_TEST_MODE" envDefault:"true"`
}

type Notifications struct {
	Telegram Telegram
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	ChatID   int64  `env:"TELEGRAM_CHAT_ID,required"`
}
