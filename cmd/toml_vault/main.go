package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"

	"github.com/reproio/toml_vault"
)

var (
	key     string
	output  string
	config  toml_vault.Config
)

func main() {
	app := cli.NewApp()

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "cryptor",
			Usage:       "cryptor",
			Value:       "simple",
			Destination: &config.Cryptor,
		},

		cli.StringFlag{
			Name:        "salt, s",
			Destination: &config.Salt,
		},

		cli.StringFlag{
			Name:        "cipher",
			Usage:       "Encrypt cipher (see. OpenSSL::Cipher.ciphers)",
			Destination: &config.Cipher,
			Value:       "aes-256-cbc",
		},
		cli.IntFlag{
			Name:        "key_len",
			Usage:       "key length of cipher",
			Destination: &config.KeyLen,
			Value:       32,
		},

		cli.StringFlag{
			Name:        "digest",
			Usage:       "Sign digest algorithm (see. OpenSSL::Digest.constants)",
			Destination: &config.Digest,
			Value:       "SHA256",
		},
		cli.IntFlag{
			Name:        "signature_key_len",
			Usage:       "key length of signature",
			Destination: &config.SignatureKeyLen,
			Value:       64,
		},

		cli.BoolFlag{
			Name:        "use_sign_passphrase",
			Destination: &config.UseSignPassphrase,
		},

		cli.StringFlag{
			Name:        "output, o",
			Usage:       "output results into `FILE`",
			Destination: &output,
		},

		cli.StringFlag{
			Name:        "aws-kms-key-id",
			Destination: &config.AwsKmsKeyID,
		},
		cli.StringFlag{
			Name:        "aws-region",
			EnvVar:      "AWS_REGION",
			Destination: &config.AwsRegion,
		},
		cli.StringFlag{
			Name:        "aws-access-key-id",
			EnvVar:      "AWS_ACCESS_KEY_ID",
			Destination: &config.AwsAccessKeyID,
		},
		cli.StringFlag{
			Name:        "aws-secret-access-key",
			EnvVar:      "AWS_SECRET_ACCESS_KEY",
			Destination: &config.AwsSecretAccessKey,
		},

		cli.StringFlag{
			Name:        "gcp-kms-resource-id",
			Destination: &config.GcpKmsResourceID,
		},
		cli.StringFlag{
			Name:        "gcp-credential-file",
			Destination: &config.GcpCredentialFile,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "encrypt",
			Flags: append(
				flags,
				cli.StringFlag{
					Name:        "key, k",
					Usage:       "key",
					Destination: &key,
				},
			),
			Action: encrypt,
		},
		{
			Name:   "decrypt",
			Flags:  flags,
			Action: decrypt,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func encrypt(c *cli.Context) error {
	return command(c, func(p *toml_vault.Placeholder, meta toml.MetaData) error {
		if err := config.GetPassphrase(); err != nil {
			return err
		}

		keys, err := toml_vault.ParseKeys(key, meta.Keys())
		if err != nil {
			return err
		}

		if err := toml_vault.Encrypt(config, p, keys); err != nil {
			return err
		}
		return nil
	})
}

func decrypt(c *cli.Context) error {
	return command(c, func(p *toml_vault.Placeholder, meta toml.MetaData) error {
		if err := config.GetPassphrase(); err != nil {
			return err
		}

		keys, err := toml_vault.ParseKeys(key, meta.Keys())
		if err != nil {
			return err
		}

		if err := toml_vault.Decrypt(config, p, keys); err != nil {
			return err
		}
		return nil
	})
}

func command(c *cli.Context, convert func(*toml_vault.Placeholder, toml.MetaData) error) error {
	r, err := getInput(c.Args().Get(0))
	if err != nil {
		return err
	}

	p := toml_vault.Placeholder{}

	meta, err := toml.DecodeReader(r, &p)
	if err != nil {
		return err
	}

	if err = convert(&p, meta); err != nil {
		return err
	}

	w, err := getOutput()
	if err != nil {
		return err
	}

	encoder := toml.NewEncoder(w)
	if err := encoder.Encode(&p); err != nil {
		return err
	}

	return nil
}

func getInput(input string) (*os.File, error) {
	if input != "" {
		return os.OpenFile(input, os.O_RDONLY, 0644)
	} else {
		return os.Stdin, nil
	}
}

func getOutput() (*os.File, error) {
	if output != "" {
		return os.Create(output)
	} else {
		return os.Stdout, nil
	}
}
