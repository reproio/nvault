package toml_vault

type Encryptor interface {
	Encrypt(interface{}) (interface{}, error)
}

func NewEncryptor(config Config) (Encryptor, error) {
	var encryptor Encryptor
	switch config.Cryptor {
	case "aws-kms", "kms":
		encryptor = &AwsCryptor{config.AwsConfig}
	case "gcp-kms":
		encryptor = &GcpCryptor{config.GcpConfig}
	case "simple":
		fallthrough
	default:
		encryptor = &SimpleCryptor{}
	}

	return encryptor, nil
}
