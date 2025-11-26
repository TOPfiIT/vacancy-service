package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (d *DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.DBName,
		d.SSLMode,
	)
}

func (s *ServerConfig) GetPortString() string {
	return fmt.Sprintf("%d", s.Port)
}

func MustLoad() *Config {
	paths := []string{
		"local.yaml",
		"config/local.yaml",
		"/app/config/local.yaml",
		"./config/local.yaml",
	}

	var cfg *Config

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		paths = append([]string{envPath}, paths...)
	}

	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil {
			cfg = &Config{}
			if err := yaml.Unmarshal(data, cfg); err != nil {
				log.Printf("Failed to parse config from %s: %v", path, err)
				continue
			}

			log.Printf("✓ Config loaded from %s", path)
			break
		}
	}

	if cfg == nil {
		log.Fatal("Config file not found in any of the searched paths")
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, _ := strconv.Atoi(port)
		cfg.Server.Port = p
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		p, _ := strconv.Atoi(dbPort)
		cfg.Database.Port = p
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.DBName = dbName
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	}
	if env := os.Getenv("SERVER_ENV"); env != "" {
		cfg.Server.Env = env
	}

	return cfg
}
