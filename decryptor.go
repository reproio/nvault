package nvault

type Decryptor interface {
	Decrypt(interface{}) (interface{}, error)
}

func NewDecryptor(config Config) (Decryptor, error) {
	var decryptor Decryptor
	switch config.Cryptor {
	case "aws-kms", "kms":
		decryptor = &AwsCryptor{config.AwsConfig}
	case "gcp-kms":
		decryptor = &GcpCryptor{config.GcpConfig}
	case "simple":
		fallthrough
	default:
		decryptor = &SimpleCryptor{}
	}

	return decryptor, nil
}
