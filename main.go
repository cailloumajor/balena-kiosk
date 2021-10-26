package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sc)

	go func() {
		for {
			select {
			case s := <-sc:
				log.Printf("caught %q signal, stopping", s)
				signal.Stop(sc)
				cancel()
			case <-ctx.Done():
				return
			}
		}
	}()
}
