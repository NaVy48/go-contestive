package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config defines json configuration structure which is loaded when program is started
type Config struct {
	Address string `json:"address"`
	BaseURL string `json:"baseUrl"`
	Title   string `json:"title"`

	// FrontEndProxy is address where not /api requests should be redirected
	FrontEndProxy string `json:"frontEndProxy"`
	// FrontEndDir sets the folder from where client app static files are served. Ignored FrontEndProxy is set
	FrontEndDir string `json:"frontEndDir"`

	JWT struct {
		Secret     string `json:"secret"`
		Expiration int    `json:"expiration"` // in seconds
	} `json:"jwt"`

	Database struct {
		PostgreSQL struct {
			Address     string `json:"address"`
			Username    string `json:"username"`
			Password    string `json:"password"`
			Database    string `json:"database"`
			SSLMode     string `json:"sslmode"`
			SSLRootCert string `json:"sslrootcert"`
		} `json:"postgresql"`
	} `json:"Database"`

	Judge struct {
		Address          string         `json:"address"`
		ProblemDir       string         `json:"problemDir"`
		JudgeCredentials map[int]string `json:"judgeCredentials"`
	} `json:"judge"`
}

// ReadConfig reads configuration json from provided file
func ReadConfig(fileName string) (*Config, error) {
	log.Printf("Reading config from file %s\n", fileName)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %v", fileName, err)
	}

	cfg := &Config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file (check file structure): %v", err)
	}

	return cfg, nil
}
