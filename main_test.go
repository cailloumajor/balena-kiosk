package main

import (
	"os"
	"path/filepath"
	"testing"
)

var badDurationEnv = map[string]string{
	"PERIOD": "badDuration",
}

func TestConfig_Init(t *testing.T) {
	tests := []struct {
		name           string
		createDotEnv   bool
		readableDotEnv bool
		setEnv         map[string]string
		wantErr        bool
	}{
		{".env file readable", true, true, nil, false},
		{".env file unreadable", true, false, nil, true},
		{"no .env file", false, true, nil, false},
		{"bad duration format", false, true, badDurationEnv, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "test_config_init-*")
			if err != nil {
				t.Fatalf("error creating temporary directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(dir); err != nil {
					t.Fatalf("failed to remove %q directory: %v", dir, err)
				}
			}()

			if err := os.Chdir(dir); err != nil {
				t.Fatalf("failed to cd to %v: %v", dir, err)
			}

			if tt.createDotEnv {
				file := filepath.Join(dir, ".env")
				if err := os.WriteFile(file, []byte(""), 0644); err != nil {
					t.Fatalf("error creating %v: %v", file, err)
				}
				if !tt.readableDotEnv {
					if err := os.Chmod(file, 0200); err != nil {
						t.Fatalf("error setting %v unreadable: %v", file, err)
					}
				}
			}

			if tt.setEnv != nil {
				for k, v := range tt.setEnv {
					if err := os.Setenv(k, v); err != nil {
						t.Fatalf("error setting %v env var: %v", k, err)
					}
				}
			}

			c := &Config{}
			if err := c.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
