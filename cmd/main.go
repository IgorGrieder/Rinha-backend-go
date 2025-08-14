package main

import (
	"fmt"
	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
)

func main() {
	cfg := config.NewConfig()
	fmt.Printf("config: %+v\n", cfg)
}
