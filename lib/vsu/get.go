package vsu

import (
	"path"
	"strconv"
	"time"

	"github.com/dollarshaveclub/vault-shared-users/internal/twofa"
	"github.com/pkg/errors"
)

// Get will look up the provided secret from Vault and return the results to the caller
// but it is up to the caller to properly format and output the information to the user
func Get(config *Config, name string) (map[string]string, error) {
	if err := checkConfig(config); err != nil {
		return nil, errors.Wrap(err, "unable to validate config")
	}

	client, err := getClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Vault client")
	}

	path := path.Join(defaultPrefix, name)

	secret, err := client.Logical().Read(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to lookup secret "+path)
	}

	info := map[string]string{}

	for key, val := range secret.Data {
		if key == TotpKey {
			decodedTotp, err := twofa.DecodeKey(secret.Data[TotpKey].(string))
			if err != nil {
				return nil, err
			}
			val = strconv.Itoa(twofa.Totp(decodedTotp, time.Now(), 6))
		}
		info[key] = val.(string)
	}

	return info, nil
}
