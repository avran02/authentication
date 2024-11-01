package config

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type YmlConfigFile struct {
	CORSConfig       `yaml:"cors"`
	CookieConfigFIle `yaml:"cookie"`
}

type CookieConfigFIle struct {
	HTTPOnly    bool   `yaml:"http_only"`
	Secure      bool   `yaml:"secure"`
	SameSite    string `yaml:"same_site"`
	Domain      string `yaml:"domain"`
	Partitioned bool   `yaml:"partitioned"`
}

type CookieConfig struct {
	HTTPOnly    bool          `yaml:"http_only"`
	Secure      bool          `yaml:"secure"`
	SameSite    http.SameSite `yaml:"same_site"`
	Domain      string        `yaml:"domain"`
	Partitioned bool          `yaml:"partitioned"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

func getYmlConfig() (CORSConfig, CookieConfig) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("can't read config.yml")
	}
	defer f.Close()
	var ymlConf YmlConfigFile
	if err := yaml.NewDecoder(f).Decode(&ymlConf); err != nil {
		f.Close()
		log.Fatal("can't decode config.yml") //nolint
	}

	var ss http.SameSite
	switch ymlConf.CookieConfigFIle.SameSite {
	case "none":
		ss = http.SameSiteNoneMode
	case "lax":
		ss = http.SameSiteLaxMode
	case "strict":
		ss = http.SameSiteStrictMode
	default:
		ss = http.SameSiteDefaultMode
	}

	return ymlConf.CORSConfig, CookieConfig{
		HTTPOnly:    ymlConf.CookieConfigFIle.HTTPOnly,
		Secure:      ymlConf.CookieConfigFIle.Secure,
		SameSite:    ss,
		Domain:      ymlConf.CookieConfigFIle.Domain,
		Partitioned: ymlConf.CookieConfigFIle.Partitioned,
	}
}
