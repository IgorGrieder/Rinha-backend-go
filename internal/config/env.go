package config

import (
	"os"
	"strconv"
)

type config struct {
	PORT       int
	REDIS_ADDR string
	REDIS_PORT int
}

func (c *config) NewConfig() *config {
	port := parseInt(getEnv("PORT", "8080"))
	reddisAddr := getEnv("REDIS_ADDR", "localhost")
	reddisPort := parseInt(getEnv("REDIS_PORT", "6639"))
	return &config{PORT: port, REDIS_ADDR: reddisAddr, REDIS_PORT: reddisPort}
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
