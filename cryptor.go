package nvault

// Cryptor ...
type Cryptor interface {
	Encryptor
	Decryptor
}

// Encryptor ...
type Encryptor interface {
	Encrypt(interface{}) (interface{}, error)
}

// Decryptor ...
type Decryptor interface {
	Decrypt(interface{}) (interface{}, error)
}

// NewCryptor ...
func NewCryptor(config *Config) Cryptor {
	var cryptor Cryptor
	switch config.Cryptor {
	case "aws-kms", "kms":
		cryptor = &AwsCryptor{config.AwsConfig}
	case "gcp-kms":
		cryptor = &GcpCryptor{config.GcpConfig}
	case "simple":
		cryptor = &SimpleCryptor{config.SimpleConfig}
	case "none":
		fallthrough
	default:
		cryptor = &NoneCryptor{config.NoneConfig}
	}

	return cryptor
}

// Config ...
type Config struct {
	NoneConfig
	SimpleConfig
	AwsConfig
	GcpConfig

	Cryptor string
}

// Option ...
type Option func(c *Config) error

// NewConfig ...
func NewConfig(cryptor string, opts ...Option) *Config {
	config := &Config{
		Cryptor: cryptor,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}
