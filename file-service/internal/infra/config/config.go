package config

import (
	"bytes"
	"log/slog"
	"os"
	"strconv"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage struct {
		Path  string
		Limit string
	}
	Server struct {
		ReadHeaderTimeout int `yaml:"read-header-timeout"`
		Port              int
	}
	Auth struct {
		PAMEnabled   bool   `yaml:"enable-pam"`
		PublicKeyURL string `yaml:"public-key-url"`
	}
}

func NewConfig(path string) *Config {
	if path == "" {
		path = "./config.yaml"
	}

	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error("could not find config yaml file with the provided path", "error", err, "path", path)
		os.Exit(1)
	}

	t := template.New("configParser").Funcs(template.FuncMap{
		"envOrKey":        envOrKey,
		"envOrKeyInt":     envOrKeyInt,
		"envOrKeyBoolean": envOrKeyBoolean,
	})

	t, err = t.Parse(string(file))
	if err != nil {
		slog.Error("error while parsing template against yaml file", "error", err)
		os.Exit(1)
	}

	var buffer bytes.Buffer
	if err = t.Execute(&buffer, nil); err != nil {
		slog.Error("error while executing template against yaml file", "error", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(buffer.Bytes(), &config); err != nil {
		slog.Error("error while reading yaml application file", "error", err)
		os.Exit(1)
	}

	return &config
}

func envOrKeyInt(envVar string, defaultValue int) (int, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue, nil
	}

	cvt, err := strconv.Atoi(value)

	if err != nil {
		slog.Error("cannot convert variable value string to int", "environmentVariable", envVar, "err", err)
		return 0, err
	}

	return cvt, nil
}

func envOrKeyBoolean(envVar string, defaultValue bool) (bool, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseBool(value)
}

func envOrKey(envVar, defaultValue string) (string, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}
