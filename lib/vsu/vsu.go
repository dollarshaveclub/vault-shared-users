package vsu

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

const (
	defaultRetries    = 5
	defaultPrefix     = "secret/shared-access" // Prefix for Vault to search for users
	defaultAuthMethod = "github"
)

var (
	defaultTimeout   = 5000 * time.Millisecond // 5 seconds
	defaultTokenPath = path.Join(os.Getenv("HOME"), ".vault-token")

	// PasswordKey is the default key for passwords stored in Vault
	PasswordKey = "password"
	// UsernameKey is the default key for usernames stored in Vault
	UsernameKey = "username"
	// TotpKey is the default key for TOTP stored in Vault
	TotpKey = "totp"
	// RecoveryCodesKey is the default key for recovery codes stored in Vault
	RecoveryCodesKey = "recovery-codes"

	defaultKeys = []string{UsernameKey, PasswordKey, TotpKey, RecoveryCodesKey}
)

// Config is the main struct for running commands
type Config struct {
	AuthMethod    string
	Login         bool
	Retries       int
	TestTokenPath string
	Timeout       time.Duration
	TokenPath     string
	VaultAddr     string
}

func checkConfig(config *Config) error {
	if config.Retries < 0 {
		config.Retries = defaultRetries
	}

	if config.Timeout < 0 {
		config.Timeout = defaultTimeout
	}

	if config.TokenPath == "" {
		config.TokenPath = defaultTokenPath
	}

	if config.AuthMethod == "" {
		config.AuthMethod = defaultAuthMethod
	}

	if config.VaultAddr == "" {
		return errors.New("Vault address was not provided or found")
	}

	return nil
}

// IsDefaultKey will look up the provided key name and check it against the default
// set of keys used.
func IsDefaultKey(key string) bool {
	for _, k := range defaultKeys {
		if k == key {
			return true
		}
	}
	return false
}

func getClient(config *Config) (*vault.Client, error) {
	client, err := vault.NewClient(&vault.Config{
		Address:    config.VaultAddr,
		MaxRetries: config.Retries,
		Timeout:    config.Timeout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to vault")
	}

	data := make(map[string]interface{})
	fileToken := []byte{}
	err = errors.New("")
	f := &os.File{}
	t := ""

	tokenPath := config.TokenPath
	if config.TestTokenPath != "" {
		tokenPath = config.TestTokenPath
	}

	authToken := ""
	for err != nil {
		f, err = os.Open(tokenPath)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}
	}
	defer f.Close()

	fileToken, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read token file at path "+tokenPath)
	}
	fileToken = bytes.TrimSpace(fileToken)

	vaultAuthPath := ""
	switch config.AuthMethod {
	case "github":
		vaultAuthPath = "auth/github/login"
		if authToken != "" {
			t = authToken
		} else {
			t = string(fileToken)
		}
		data["token"] = t
	default:
		return nil, fmt.Errorf("auth method %s not implemented", config.AuthMethod)
	}

	token := ""
	if config.Login {
		secret, err := client.Logical().Write(vaultAuthPath, data)
		if err != nil {
			return nil, errors.Wrap(err, "unable to login to vault on "+vaultAuthPath)
		}

		token, err = secret.TokenID()
		if err != nil {
			return nil, errors.Wrap(err, "unable to lookup token")
		}
	} else {
		token = string(fileToken)
	}

	client.SetToken(token)

	return client, nil
}

// GetKeys will provide the keys for a given map back to the user as a list of strings.
// Note that the default keys will be presented first in the order defined by the
// "defaultKeys" variable and additional keys are added in a random order.
func GetKeys(data map[string]string) []string {
	keys := defaultKeys
	if data != nil && len(data) > len(defaultKeys) {
		for name := range data {
			if !IsDefaultKey(name) {
				keys = append(keys, name)
			}
		}
	}

	return keys
}
