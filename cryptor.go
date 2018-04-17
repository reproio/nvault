package nvault

type Cryptor interface {
	Encryptor
	Decryptor
}

type Encryptor interface {
	Encrypt(interface{}) (interface{}, error)
}

type Decryptor interface {
	Decrypt(interface{}) (interface{}, error)
}

type Config struct {
	SimpleConfig
	AwsConfig
	GcpConfig

	Cryptor string
}

func NewCryptor(config *Config) Cryptor {
	var cryptor Cryptor
	switch config.Cryptor {
	case "aws-kms", "kms":
		cryptor = &AwsCryptor{config.AwsConfig}
	case "gcp-kms":
		cryptor = &GcpCryptor{config.GcpConfig}
	case "simple":
		fallthrough
	default:
		cryptor = &SimpleCryptor{config.SimpleConfig}
	}

	return cryptor
}
