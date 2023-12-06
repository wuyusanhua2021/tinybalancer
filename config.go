package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ascii = `
_________  ___  ________       ___    ___ ________  ________  ___       ________  ________   ________  _______   ________     
|\___   ___\\  \|\   ___  \    |\  \  /  /|\   __  \|\   __  \|\  \     |\   __  \|\   ___  \|\   ____\|\  ___ \ |\   __  \    
\|___ \  \_\ \  \ \  \\ \  \   \ \  \/  / | \  \|\ /\ \  \|\  \ \  \    \ \  \|\  \ \  \\ \  \ \  \___|\ \   __/|\ \  \|\  \   
     \ \  \ \ \  \ \  \\ \  \   \ \    / / \ \   __  \ \   __  \ \  \    \ \   __  \ \  \\ \  \ \  \    \ \  \_|/_\ \   _  _\  
      \ \  \ \ \  \ \  \\ \  \   \/  /  /   \ \  \|\  \ \  \ \  \ \  \____\ \  \ \  \ \  \\ \  \ \  \____\ \  \_|\ \ \  \\  \| 
       \ \__\ \ \__\ \__\\ \__\__/  / /      \ \_______\ \__\ \__\ \_______\ \__\ \__\ \__\\ \__\ \_______\ \_______\ \__\\ _\ 
        \|__|  \|__|\|__| \|__|\___/ /        \|_______|\|__|\|__|\|_______|\|__|\|__|\|__| \|__|\|_______|\|_______|\|__|\|__|
                              \|___|/                                                                                                                                                                                                  
`

type Config struct {
	SSLCertificateKey   string      `yaml:"ssl_certificate_key"`
	Location            []*Location `yaml:"location"`
	Schema              string      `yaml:"schema"`
	Port                int         `yaml:"port"`
	SSLCertificate      string      `yaml:"ssl_certificate"`
	HealthCheck         bool        `yaml:"health_check"`
	HealthCheckInterval uint        `yaml:"health_check_interval"`
	MaxAllowed          uint        `yaml:"max_allowed"`
}

type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

func ReadConfig(fileName string) (*Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Print() {
	fmt.Printf("%s\nSchema: %s\nPort: %d\nHealth Check: %v\nLocation:\n",
		ascii,
		c.Schema,
		c.Port,
		c.HealthCheck,
	)

	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxy Pass: %s\n\tMode: %s\n\n",
			l.Pattern,
			l.ProxyPass,
			l.BalanceMode,
		)
	}
}

func (c *Config) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("the schema \"%s\" not supported", c.Schema)
	}
	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}
	if c.Schema == "https" && (len(c.SSLCertificate) == 0 || len(c.SSLCertificateKey) == 0) {
		return errors.New("the https proxy requires ssl_certificate and ssl_certificate_key")
	}
	if c.HealthCheckInterval < 1 {
		return errors.New("health_check_interval must be greater than zero")
	}
	return nil
}
