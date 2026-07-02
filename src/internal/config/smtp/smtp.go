package smtp

type SMTPConfig struct {
	Host        string `yaml:"host" env:"SMTP_HOST" env-default:""`
	Port        int    `yaml:"port" env:"SMTP_PORT" env-default:"587"`
	Username    string `yaml:"username" env:"SMTP_USERNAME" env-default:""`
	Password    string `yaml:"password" env:"SMTP_PASSWORD" env-default:""`
	SenderName  string `yaml:"sender_name" env:"SMTP_SENDER_NAME" env-default:"GreenMart"`
	SenderEmail string `yaml:"sender_email" env:"SMTP_SENDER_EMAIL" env-default:""`
}
