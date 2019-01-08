package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/dollarshaveclub/vault-shared-users/lib/vsu"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [flags] <name>",
	Short: "Add or update a set of credentials",
	Long:  `Add or update a set of credentials`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		existing, err := vsu.GetExisting(config, name)
		if err != nil {
			cliError("unable to look for existing secret", err)
		}

		creds, err := credInput(existing)
		if err != nil {
			cliError("unable to get creds from user", err)
		}

		err = vsu.Add(config, name, creds)
		if err != nil {
			cliError("unable to update credentials", err)
		}

	},
}

func credInput(existing map[string]string) (map[string]interface{}, error) {
	var err error
	buf := bufio.NewReader(os.Stdin)

	creds := map[string]interface{}{}

	keys := vsu.GetKeys(existing)

	for _, key := range keys {
		fmt.Print(getQuestion(key, existing))
		value := ""
		if key == vsu.PasswordKey {
			v, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return nil, errors.Wrap(err, "unable to read password value")
			}
			value = string(v)
			fmt.Println()
		} else {
			value, err = buf.ReadString('\n')
			if err != nil {
				return nil, errors.Wrap(err, "unable to read value for "+key)
			}

		}
		value = strings.TrimSpace(value)
		if value == "" {
			if d, ok := existing[key]; ok {
				value = d
			}
		} else if strings.HasPrefix(value, "@") {
			value, err = getValueFromFile(value[1:])
			if err != nil {
				return nil, errors.Wrap(err, "unable to read user provided file")
			}
		}

		creds[key] = value
	}

	fmt.Println("If you would like to add any other information, do so in the next section")
	fmt.Println("When done, enter and empty key name to finish")
	for {
		fmt.Print("key: ")
		key, err := buf.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "unable to read custom key")
		}
		key = strings.TrimSpace(key)

		if key == "" {
			break
		}

		fmt.Print("value: ")
		value, err := buf.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "unable to read custom value")
		}
		value = strings.TrimSpace(value)

		creds[key] = value
	}

	return creds, nil
}

func getQuestion(key string, existing map[string]string) string {
	str := fmt.Sprintf("Please enter a value for %s", key)
	if existing != nil {
		last := 1
		if len(existing[key]) > 4 {
			last = 4
		}
		if d, ok := existing[key]; ok {
			if key == vsu.PasswordKey {
				d = d[0:last] + "..."
			}
			str += " (current: " + d + ")"
		}
	}
	return str + ": "
}

func getValueFromFile(filename string) (string, error) {
	filename, err := homedir.Expand(os.ExpandEnv(filename))
	if err != nil {
		return "", errors.Wrap(err, "unable to expand filename")
	}

	f, err := os.Open(filename)
	if err != nil {
		return "", errors.Wrap(err, "unable to open file "+filename)
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "unable to read file "+filename)
	}
	return string(contents), nil
}
