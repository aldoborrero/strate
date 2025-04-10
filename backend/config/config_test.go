package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		validate   func(*testing.T, *Config)
		wantErr    string
	}{
		{
			name:       "basic_config",
			configPath: "testdata/basic.yaml",
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 4000, cfg.Server.Port)
				assert.Equal(t, "https://strate.example.com", cfg.Server.BaseURL)
				assert.Equal(t, "sqlite", cfg.Database.Type)
				assert.Equal(t, "/tmp/strate.db", cfg.Database.Sqlite.Path)
				assert.Equal(t, "github", cfg.CI.Provider)
				assert.Equal(t, "github.enterprise.com", cfg.CI.GitHub.Hostname)
				assert.Equal(t, "test-secret", cfg.CI.GitHub.WebhookSecret)
				assert.Equal(t, "github-token", cfg.CI.GitHub.Token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the global config
			BackendConfig = nil

			cfg, err := LoadConfig(tt.configPath)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, cfg, BackendConfig)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		envVars    map[string]string
		validate   func(*testing.T, *Config)
	}{
		{
			name:       "basic_env_vars",
			configPath: "testdata/basic.yaml",
			envVars: map[string]string{
				"STRATE_SERVER_PORT":            "5000",
				"STRATE_SERVER_BASE_URL":        "https://strate-env.example.com",
				"STRATE_DATABASE_TYPE":          "postgres",
				"STRATE_DATABASE_POSTGRES_HOST": "postgres-env-host",
				"STRATE_DATABASE_POSTGRES_PORT": "5433",
				"STRATE_DATABASE_POSTGRES_NAME": "strate-env-db",
				"STRATE_DATABASE_POSTGRES_USER": "postgres-env-user",
				"STRATE_DATABASE_POSTGRES_PASS": "postgres-env-password",
				"STRATE_CI_PROVIDER":            "github",
				"STRATE_CI_GITHUB_HOSTNAME":     "github-env.enterprise.com",
				"STRATE_AUTH_JWT_SECRET":        "env-jwt-secret",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, 5000, cfg.Server.Port)
				assert.Equal(t, "https://strate-env.example.com", cfg.Server.BaseURL)
				assert.Equal(t, "postgres", cfg.Database.Type)
				assert.Equal(t, "postgres-env-host", cfg.Database.Postgres.Host)
				assert.Equal(t, 5433, cfg.Database.Postgres.Port)
				assert.Equal(t, "strate-env-db", cfg.Database.Postgres.Name)
				assert.Equal(t, "postgres-env-user", cfg.Database.Postgres.User)
				assert.Equal(t, "postgres-env-password", cfg.Database.Postgres.Pass)
				assert.Equal(t, "github", cfg.CI.Provider)
				assert.Equal(t, "github-env.enterprise.com", cfg.CI.GitHub.Hostname)
				assert.Equal(t, "env-jwt-secret", cfg.Auth.JWTSecret)
			},
		},
		{
			name:       "feature_flags",
			configPath: "testdata/basic.yaml",
			envVars: map[string]string{
				"STRATE_FEATURES_USER_SERVICE_ENABLED":                "false",
				"STRATE_FEATURES_LIMIT_MAX_PROJECTS_TO_FILES_CHANGED": "true",
				"STRATE_FEATURES_INTERNAL_USERS_ENABLED":              "true",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.False(t, cfg.Features.UserServiceEnabled)
				assert.True(t, cfg.Features.LimitMaxProjectsToFilesChanged)
				assert.True(t, cfg.Features.InternalUsersEnabled)
			},
		},
		{
			name:       "log_settings",
			configPath: "testdata/basic.yaml",
			envVars: map[string]string{
				"STRATE_LOG_LEVEL":  "debug",
				"STRATE_LOG_FORMAT": "json",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "json", cfg.Log.Format)
				assert.Equal(t, int(cfg.Log.Level), int(slog.LevelDebug))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the global config
			BackendConfig = nil

			// Set environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := LoadConfig(tt.configPath)
			require.NoError(t, err)
			assert.NotNil(t, cfg)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}

			// Verify global config is set
			assert.Equal(t, cfg, BackendConfig)
		})
	}
}
