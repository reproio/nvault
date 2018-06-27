package cmd

import (
	"io"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/reproio/nvault"
)

var (
	key    string
	output string
	config nvault.Config
)

// Converter ...
type Converter func(input io.Reader, output io.Writer, cryptor Cryptor) error

// Cryptor ...
type Cryptor func(*nvault.Placeholder) error

// Run ...
func Run(converter Converter) {
	app := cli.NewApp()
	app.Commands = subcommands(converter)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// RunAll ...
func RunAll() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:        "json",
			Subcommands: subcommands(JSONConverter),
		},
		{
			Name:        "toml",
			Subcommands: subcommands(TomlConverter),
		},
		{
			Name:        "yaml",
			Subcommands: subcommands(YamlConverter),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func subcommands(converter Converter) []cli.Command {
	return []cli.Command{
		{
			Name:   "encrypt",
			Flags:  flags(),
			Action: command(converter, encryptor),
		},
		{
			Name:   "decrypt",
			Flags:  flags(),
			Action: command(converter, decryptor),
		},
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
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
			Destination: &config.AwsRegion,
		},
		cli.StringFlag{
			Name:        "aws-access-key-id",
			Destination: &config.AwsAccessKeyID,
		},
		cli.StringFlag{
			Name:        "aws-secret-access-key",
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
}

func command(converter Converter, cryptor Cryptor) func(*cli.Context) error {
	return func(c *cli.Context) (err error) {
		if err = GetPassphrase(&config); err != nil {
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

func encryptor(p *nvault.Placeholder) (err error) {
	var paths []nvault.Path
	if key != "" {
		paths, err = ParseKeys(key)
		if err != nil {
			return
		}
	} else {
		for _, path := range p.Paths() {
			paths = append(paths, path.AddRoot(nvault.PathFragment{Type: "string", Fragment: "$"}))
		}
	}

	if err = nvault.Encrypt(p, &config, paths...); err != nil {
		return
	}
	return
}

func decryptor(p *nvault.Placeholder) (err error) {
	var paths []nvault.Path
	if key != "" {
		paths, err = ParseKeys(key)
		if err != nil {
			return
		}
	} else {
		for _, path := range p.Paths() {
			paths = append(paths, path.AddRoot(nvault.PathFragment{Type: "string", Fragment: "$"}))
		}
	}

	if err = nvault.Decrypt(p, &config, paths...); err != nil {
		return
	}
	return
}

func getReader(input string) (*os.File, error) {
	if input != "" {
		return os.OpenFile(input, os.O_RDONLY, 0644)
	}
	return os.Stdin, nil
}

func getWriter() (*os.File, error) {
	if output != "" {
		return os.Create(output)
	}
	return os.Stdout, nil
}
