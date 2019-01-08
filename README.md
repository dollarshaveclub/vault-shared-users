# Vault Shared Users

Vault Shared Users provide us with a mechanism to share credentials for users that are not easily shared such as bot accounts for GitHub.

One of the main features of this tool is that it will also store and generate 2-factor authentication (2FA) codes for these accounts while still allowing multiple users access. The 2FA code seed is stored and used by this tool to produce time-based one time passwords (TOTP) codes.

## Vault

You should have access to Vault before using this tool. You should also login to Vault if you use a method other than GitHub based authentication.

## Listing available credentials

The `list` subcommand will provide you with the names of the avialable credentials to look up.

## Getting credentials

The `get` subcommand will provide you with the crendial information stored it Vault. A code for the current time period will be generated for 2FA automatically.

If you just need a 2FA code, you can use the `--2fa-only` flag and the other information will be skipped.

## Adding a new set of credentials

The `add` subcommand can be used to add or update a set of credentials. The command will ask for the following information my default:

* `username`
* `password`
* `totp`
* `recovery-codes`

For `totp`, the seed to generate the codes should be provided by the site. If all you see is a QR code, there will usually be a link to allow you to type in the code into your device. You will want to copy that code and add it to the credentials as it is provided.

For `recovery-codes`, the site will provide us with a set of one-use codes when we need to reset our MFA device. These codes will usually be provided in a file and downloaded. If you have a file, you can prefix the filename with an "@" sign to specify that the codes should be consumed from a file. For example, if your file is located at `/tmp/recovery-codes.txt` you would specify the codes with `@/tmp/recovery-codes.txt`.

After the default information is added, the command will prompt for any additional information as well. These additional pieces of informatino can include, but are not limited to:

* email address used for sign-up
* security verification questions such as "What's your Mother's maiden name?" or "What high school did you go to?" When completing these questions when setting up the account you should add in random generated answers.
* anything else that might be helpful to know with this account.
