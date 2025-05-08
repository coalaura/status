package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	StatusPage string

	HTTPTimeout    int
	HTTPRetryDelay int
	HTTPRetryCount int

	EmailTo string

	SMTPHost     string
	SMTPPort     int
	SMTPPassword string
	SMTPUser     string

	TemplateFavicon     string
	TemplateBanner      string
	TemplateURL         string
	TemplateTitle       string
	TemplateDescription string
	TemplateCopyright   string
}

func ReadMainConfig() (*Config, error) {
	data, err := os.ReadFile(".env")
	if err != nil {
		return nil, err
	}

	env, err := godotenv.Unmarshal(string(data))
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(env["SMTP_PORT"])

	timeout, _ := strconv.Atoi(env["HTTP_TIMEOUT"])
	delay, _ := strconv.Atoi(env["HTTP_RETRY_DELAY"])
	retry, _ := strconv.Atoi(env["HTTP_RETRY_COUNT"])

	if timeout == 0 {
		timeout = 5
	}

	if delay == 0 {
		delay = 10
	}

	return &Config{
		StatusPage: env["STATUS_PAGE"],

		HTTPTimeout:    timeout,
		HTTPRetryDelay: delay,
		HTTPRetryCount: retry,

		EmailTo:      env["EMAIL_TO"],
		SMTPHost:     env["SMTP_HOST"],
		SMTPPort:     port,
		SMTPUser:     env["SMTP_USER"],
		SMTPPassword: env["SMTP_PASSWORD"],

		TemplateFavicon:     env["TEMPLATE_FAVICON"],
		TemplateBanner:      env["TEMPLATE_BANNER"],
		TemplateURL:         env["TEMPLATE_URL"],
		TemplateTitle:       env["TEMPLATE_TITLE"],
		TemplateDescription: env["TEMPLATE_DESCRIPTION"],
		TemplateCopyright:   env["TEMPLATE_COPYRIGHT"],
	}, nil
}
