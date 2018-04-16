package cmd

import (
	"io"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/reproio/toml_vault"
)

type Converter func(input io.Reader, output io.Writer, cryptor Cryptor) error
type Cryptor func(*toml_vault.Placeholder) error

var (
	key    string
	output string
	config toml_vault.Config
)

func Run(converter Converter) {
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
			Name:        "key, k",
			Usage:       "key",
			Destination: &key,
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
			Name:   "encrypt",
			Flags:  flags,
			Action: command(converter, encryptor),
		},
		{
			Name:   "decrypt",
			Flags:  flags,
			Action: command(converter, decryptor),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func command(converter Converter, cryptor Cryptor) func(*cli.Context) error {
	return func(c *cli.Context) (err error) {
		if err = config.GetPassphrase(); err != nil {
			return err
		}

		input, err := getReader(c.Args().Get(0))
		if err != nil {
			return err
		}

		output, err := getWriter()
		if err != nil {
			return err
		}

		return converter(input, output, cryptor)
	}
}

func encryptor(p *toml_vault.Placeholder) (err error) {
	var paths []toml_vault.Path
	if key != "" {
		paths, err = ParseKeys(key)
		if err != nil {
			return
		}
	} else {
		for _, path := range p.Paths() {
			paths = append(paths, path.AddRoot(toml_vault.PathFragment{"string", "$"}))
		}
	}

	if err = toml_vault.Encrypt(config, p, paths); err != nil {
		return
	}
	return
}

func decryptor(p *toml_vault.Placeholder) (err error) {
	var paths []toml_vault.Path
	if key != "" {
		paths, err = ParseKeys(key)
		if err != nil {
			return
		}
	} else {
		for _, path := range p.Paths() {
			paths = append(paths, path.AddRoot(toml_vault.PathFragment{"string", "$"}))
		}
	}

	if err = toml_vault.Decrypt(config, p, paths); err != nil {
		return
	}
	return
}

func getReader(input string) (*os.File, error) {
	if input != "" {
		return os.OpenFile(input, os.O_RDONLY, 0644)
	} else {
		return os.Stdin, nil
	}
}

func getWriter() (*os.File, error) {
	if output != "" {
		return os.Create(output)
	} else {
		return os.Stdout, nil
	}
}
