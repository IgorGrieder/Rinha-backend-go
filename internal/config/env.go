package config

import (
	"os"
	"strconv"
)

type Config struct {
	PORT          int
	REDIS_ADDR    string
	DEFAULT_ADDR  string
	FALLBACK_ADDR string
	REDIS_PORT    int
	QUEUE         string
}

func NewConfig() *Config {
	port := parseInt(getEnv("PORT", "8080"))
	reddisAddr := getEnv("REDIS_ADDR", "localhost")
	reddisPort := parseInt(getEnv("REDIS_PORT", "6639"))
	queue := getEnv("QUEUE", "payment-processor-queue")
	fallbackAddr := getEnv("FALLBACK_ADDR", "")
	defaultAddr := getEnv("DEFAULT_ADDR", "")
	return &Config{
		PORT:          port,
		REDIS_ADDR:    reddisAddr,
		REDIS_PORT:    reddisPort,
		QUEUE:         queue,
		DEFAULT_ADDR:  defaultAddr,
		FALLBACK_ADDR: fallbackAddr,
	}
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
