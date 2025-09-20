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
	HASH_DEFAULT  string
	HASH_FALLBACK string
	WORKERS       int
}

func NewConfig() *Config {
	port := parseInt(getEnv("PORT", "8080"))
	reddisAddr := getEnv("REDIS_ADDR", "redis")
	reddisPort := parseInt(getEnv("REDIS_PORT", "6379"))
	workers := parseInt(getEnv("WORKERS", "2"))
	queue := getEnv("QUEUE", "payment-processor-queue")
	fallbackAddr := getEnv("PAYMENT_PROCESSOR_FALLBACK", "")
	defaultAddr := getEnv("PAYMENT_PROCESSOR_DEFAULT", "")
	hashDef := getEnv("HASH_DEFAULT", "")
	hashFall := getEnv("HASH_FALLBACK", "")
	return &Config{
		PORT:          port,
		REDIS_ADDR:    reddisAddr,
		REDIS_PORT:    reddisPort,
		QUEUE:         queue,
		DEFAULT_ADDR:  defaultAddr,
		FALLBACK_ADDR: fallbackAddr,
		WORKERS:       workers,
		HASH_DEFAULT:  hashDef,
		HASH_FALLBACK: hashFall,
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
