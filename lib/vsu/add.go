package vsu

import (
	"path"

	"github.com/pkg/errors"
)

// Add will take in a set of credentials and add or update them inside of Vault. It
// is up to the caller to check for an existing set of credentials and to merge those
// with the new credential updates before calling add. Otherwise, add will clobber
// any old values that were not provided
func Add(config *Config, name string, creds map[string]interface{}) error {
	if err := checkConfig(config); err != nil {
		return errors.Wrap(err, "unable to validate config")
	}

	client, err := getClient(config)
	if err != nil {
		return errors.Wrap(err, "unable to create Vault client")
	}

	vPath := path.Join(defaultPrefix, name)

	_, err = client.Logical().Write(vPath, creds)
	if err != nil {
		return errors.Wrap(err, "unable to update credentials")
	}

	return nil
}

// GetExisting will look for an existing set of credentials with the provided
// name and will return the data for usage by user fronting programs
func GetExisting(config *Config, name string) (map[string]string, error) {
	if err := checkConfig(config); err != nil {
		return nil, errors.Wrap(err, "unable to validate config")
	}

	client, err := getClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Vault client")
	}

	vPath := path.Join(defaultPrefix, name)

	existing, err := client.Logical().Read(vPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to look for existing secret")
	}

	info := map[string]string{}

	if existing != nil && existing.Data != nil {
		for key, val := range existing.Data {
			info[key] = val.(string)
		}
	}

	return info, nil
}
