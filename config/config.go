package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		MySQL struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"mysql"`
		Redis struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"db"`

	Connection struct {
		HTTP struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			ReadTimeout  int    `yaml:"readTimeout"`
			WriteTimeout int    `yaml:"writeTimeout"`
			IdleTimeout  int    `yaml:"idleTimeout"`
		} `yaml:"http"`
	} `yaml:"connection"`

	Notification struct {
		SMTPHost    string `yaml:"smtpHost"`
		SMTPPort    int    `yaml:"smtpPort"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		SenderEmail string `yaml:"senderEmail"`
	} `yaml:"notification"`
}

// LoadConfig reads the YAML config from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
