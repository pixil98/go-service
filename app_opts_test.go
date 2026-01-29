package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCmdLineOpts(t *testing.T) {
	opts, err := newCmdLineOpts("test")
	if err != nil {
		t.Fatalf("newCmdLineOpts() error = %v", err)
	}

	if opts.fs.Lookup("config") == nil {
		t.Error("config flag not registered")
	}
	if opts.fs.Lookup("loglevel") == nil {
		t.Error("loglevel flag not registered")
	}
}

func TestCmdLineOpts_Parse(t *testing.T) {
	tests := map[string]struct {
		args         []string
		wantConfig   string
		wantLoglevel string
		wantErr      bool
	}{
		"default values": {
			args:         []string{"cmd"},
			wantConfig:   "",
			wantLoglevel: "info",
			wantErr:      false,
		},
		"config flag": {
			args:         []string{"cmd", "-config", "/path/to/config.json"},
			wantConfig:   "/path/to/config.json",
			wantLoglevel: "info",
			wantErr:      false,
		},
		"loglevel flag": {
			args:         []string{"cmd", "-loglevel", "debug"},
			wantConfig:   "",
			wantLoglevel: "debug",
			wantErr:      false,
		},
		"both flags": {
			args:         []string{"cmd", "-config", "/path/to/config.json", "-loglevel", "warn"},
			wantConfig:   "/path/to/config.json",
			wantLoglevel: "warn",
			wantErr:      false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opts, _ := newCmdLineOpts("test")
			err := opts.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if opts.config != tt.wantConfig {
				t.Errorf("config = %v, want %v", opts.config, tt.wantConfig)
			}
			if opts.loglevel != tt.wantLoglevel {
				t.Errorf("loglevel = %v, want %v", opts.loglevel, tt.wantLoglevel)
			}
		})
	}
}

func TestCmdLineOpts_Logger(t *testing.T) {
	tests := map[string]struct {
		loglevel string
		wantErr  bool
	}{
		"debug level": {
			loglevel: "debug",
			wantErr:  false,
		},
		"info level": {
			loglevel: "info",
			wantErr:  false,
		},
		"warn level": {
			loglevel: "warn",
			wantErr:  false,
		},
		"error level": {
			loglevel: "error",
			wantErr:  false,
		},
		"invalid level": {
			loglevel: "invalid",
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opts := &cmdLineOpts{loglevel: tt.loglevel, logformat: "json"}
			logger, err := opts.Logger()
			if (err != nil) != tt.wantErr {
				t.Errorf("Logger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("Logger() returned nil logger")
			}
		})
	}
}

func TestCmdLineOpts_Config(t *testing.T) {
	type testConfig struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := map[string]struct {
		configContent string
		setupConfig   bool
		wantErr       bool
		wantName      string
		wantValue     int
	}{
		"valid config": {
			configContent: `{"name": "test", "value": 42}`,
			setupConfig:   true,
			wantErr:       false,
			wantName:      "test",
			wantValue:     42,
		},
		"empty config path": {
			configContent: "",
			setupConfig:   false,
			wantErr:       true,
		},
		"invalid json": {
			configContent: `not json`,
			setupConfig:   true,
			wantErr:       true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opts := &cmdLineOpts{}

			if tt.setupConfig {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "config.json")
				if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
					t.Fatalf("failed to write test config: %v", err)
				}
				opts.config = configPath
			}

			var cfg testConfig
			err := opts.Config(&cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg.Name != tt.wantName {
					t.Errorf("config.Name = %v, want %v", cfg.Name, tt.wantName)
				}
				if cfg.Value != tt.wantValue {
					t.Errorf("config.Value = %v, want %v", cfg.Value, tt.wantValue)
				}
			}
		})
	}
}
