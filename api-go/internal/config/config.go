package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds env for server, DB, Redis, and external APIs.
type Config struct {
	Port        string
	DatabaseURL string
	BaseURL     string
	Production  bool

	AIModerationEnabled      bool
	AIModerationPollInterval int
	AIModerationBatchSize    int
	YandexFolderID           string
	YandexAPIKey             string

	RedisAddr     string
	RedisPassword string

	MTSBearer       string
	NotisendAPIKey  string
	NotisendProject string

	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	SMTPFrom        string
	SMTPSecure      bool
	SMTPTLSInsecure bool

	S3Endpoint   string
	S3Region     string
	S3Bucket     string
	S3AccessKey  string
	S3SecretKey  string
	S3PublicBase string

	TinkoffTerminalKey string
	TinkoffSecretKey   string

	DadataAPIKey string

	CDEKClientID     string
	CDEKClientSecret string
	CDEKAPIBase      string

	DealPlatformFeePercent int
	DealPayoutDelayDays    int
	DealAutoCompleteDays   int
}

func Load() Config {
	p := os.Getenv("PORT")
	if p == "" {
		p = "3000"
	}
	aiEnabled := strings.ToLower(strings.TrimSpace(os.Getenv("AI_MODERATION_ENABLED")))
	aiModerationEnabled := aiEnabled != "false" && aiEnabled != "0" && aiEnabled != "no"
	aiPollInterval := 30
	if v := strings.TrimSpace(os.Getenv("AI_MODERATION_POLL_INTERVAL_SECONDS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			aiPollInterval = n
		}
	}
	aiBatchSize := 5
	if v := strings.TrimSpace(os.Getenv("AI_MODERATION_BATCH_SIZE")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			aiBatchSize = n
		}
	}
	base := os.Getenv("BASE_URL")
	if base == "" {
		base = "http://localhost:3000"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis"
	}
	smtpPort := 587
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			smtpPort = n
		}
	}
	sec := strings.ToLower(strings.TrimSpace(os.Getenv("SMTP_SECURE")))
	smtpSecure := sec == "true" || sec == "1" || sec == "yes"
	smtpTLSInsecure := strings.ToLower(strings.TrimSpace(os.Getenv("SMTP_TLS_INSECURE"))) == "true" ||
		strings.TrimSpace(os.Getenv("SMTP_TLS_INSECURE")) == "1"
	feePercent := 0
	if v := strings.TrimSpace(os.Getenv("DEAL_PLATFORM_FEE_PERCENT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			feePercent = n
		}
	}
	payoutDelay := 1
	if v := strings.TrimSpace(os.Getenv("DEAL_PAYOUT_DELAY_DAYS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			payoutDelay = n
		}
	}
	autoComplete := 14
	if v := strings.TrimSpace(os.Getenv("DEAL_AUTO_COMPLETE_DAYS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			autoComplete = n
		}
	}
	cdekBase := strings.TrimSpace(os.Getenv("CDEK_API_BASE"))
	if cdekBase == "" {
		cdekBase = "https://api.cdek.ru/v2"
	}

	return Config{
		Port:        p,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		BaseURL:     strings.TrimRight(base, "/"),
		Production:  strings.ToLower(os.Getenv("NODE_ENV")) == "production",

		AIModerationEnabled:      aiModerationEnabled,
		AIModerationPollInterval: aiPollInterval,
		AIModerationBatchSize:    aiBatchSize,
		YandexFolderID:           strings.TrimSpace(os.Getenv("YANDEX_FOLDER_ID")),
		YandexAPIKey:             strings.TrimSpace(os.Getenv("YANDEX_API_KEY")),

		RedisAddr:     redisHost + ":" + redisPort,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		MTSBearer:       os.Getenv("MTS_TOKEN"),
		NotisendAPIKey:  os.Getenv("NOTISEND_API_KEY"),
		NotisendProject: firstNonEmpty(os.Getenv("NOTISEND_PROJECT"), "torgui_sam"),

		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        smtpPort,
		SMTPUser:        os.Getenv("SMTP_USER"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        os.Getenv("SMTP_FROM"),
		SMTPSecure:      smtpSecure,
		SMTPTLSInsecure: smtpTLSInsecure,

		S3Endpoint:   firstNonEmpty(os.Getenv("S3_ENDPOINT"), "https://s3.ru1.storage.beget.cloud"),
		S3Region:     firstNonEmpty(os.Getenv("S3_REGION"), "ru1"),
		S3Bucket:     os.Getenv("S3_BUCKET_NAME"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:  os.Getenv("S3_SECRET_KEY"),
		S3PublicBase: os.Getenv("S3_PUBLIC_BASE"),

		TinkoffTerminalKey: strings.TrimSpace(os.Getenv("TINKOFF_TERMINAL_KEY")),
		TinkoffSecretKey:   strings.TrimSpace(os.Getenv("TINKOFF_SECRET_KEY")),

		DadataAPIKey: strings.TrimSpace(os.Getenv("DADATA_API_KEY")),

		CDEKClientID:     strings.TrimSpace(os.Getenv("CDEK_CLIENT_ID")),
		CDEKClientSecret: strings.TrimSpace(os.Getenv("CDEK_CLIENT_SECRET")),
		CDEKAPIBase:      cdekBase,

		DealPlatformFeePercent: feePercent,
		DealPayoutDelayDays:    payoutDelay,
		DealAutoCompleteDays:   autoComplete,
	}
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return b
}
