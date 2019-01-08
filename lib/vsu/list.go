package vsu

import (
	"github.com/pkg/errors"
)

// List will provide the available credentials that the user can lookup.
func List(config *Config) ([]string, error) {
	if err := checkConfig(config); err != nil {
		return nil, errors.Wrap(err, "unable to validate config")
	}

	client, err := getClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Vault client")
	}

	creds, err := client.Logical().List(defaultPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list available credentials from Vault")
	}

	list := []string{}

	if data, ok := creds.Data["keys"]; ok {
		for _, name := range data.([]interface{}) {
			list = append(list, name.(string))
		}
	}

	return list, nil
}
