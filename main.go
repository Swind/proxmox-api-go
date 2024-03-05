package main

import (
	"crypto/tls"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/Telmate/proxmox-api-go/cli"
	_ "github.com/Telmate/proxmox-api-go/cli/command/commands"
	"github.com/Telmate/proxmox-api-go/proxmox"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APIURL      string
	HTTPHeaders string
	User        string
	Password    string
	OTP         string
	NewCLI      bool
}

func loadAppConfig() AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	return AppConfig{
		APIURL:      os.Getenv("PM_API_URL"),
		HTTPHeaders: os.Getenv("PM_HTTP_HEADERS"),
		User:        os.Getenv("PM_USER"),
		Password:    os.Getenv("PM_PASS"),
		OTP:         os.Getenv("PM_OTP"),
		NewCLI:      os.Getenv("NEW_CLI") == "true",
	}
}

func initializeProxmoxClient(config AppConfig, secure bool, proxyURL string, taskTimeout int) (*proxmox.Client, error) {
	tlsconf := &tls.Config{InsecureSkipVerify: !secure}
	if secure {
		tlsconf = nil
	}

	client, err := proxmox.NewClient(config.APIURL, nil, config.HTTPHeaders, tlsconf, proxyURL, taskTimeout)
	if err != nil {
		return nil, err
	}

	if userRequiresAPIToken(config.User) {
		client.SetAPIToken(config.User, config.Password)
		_, err := client.GetVersion()
		if err != nil {
			return nil, err
		}
	} else {
		err = client.Login(config.User, config.Password, config.OTP)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func main() {
	loadAppConfig()
	err := cli.Execute()
	if err != nil {
		failError(err)
	}
	os.Exit(0)
}

var rxUserRequiresToken = regexp.MustCompile("[a-z0-9]+@[a-z0-9]+![a-z0-9]+")

func userRequiresAPIToken(userID string) bool {
	return rxUserRequiresToken.MatchString(userID)
}

// GetConfig get config from file
func GetConfig(configFile string) (configSource []byte) {
	var err error
	if configFile != "" {
		configSource, err = os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		configSource, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func failError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
