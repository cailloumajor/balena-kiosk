package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	Period time.Duration `envconfig:"default=1m"`
}

func (c *Config) Init() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	if err := envconfig.Init(c); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	return nil
}

func main() {
	c := &Config{}
	if err := c.Init(); err != nil {
		log.Fatal(err)
	}
}
