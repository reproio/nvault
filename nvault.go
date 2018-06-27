package nvault

// Encrypt ...
func Encrypt(p *Placeholder, config *Config, paths ...Path) error {
	cryptor := NewCryptor(config)

	for _, path := range p.Matches(paths) {
		value, err := p.Get(path)
		if err != nil {
			return err
		}
		value, err = cryptor.Encrypt(value)
		if err != nil {
			return err
		}
		p.Set(path, value)
	}
	return nil
}

// Decrypt ...
func Decrypt(p *Placeholder, config *Config, paths ...Path) error {
	cryptor := NewCryptor(config)

	for _, path := range p.Matches(paths) {
		value, err := p.Get(path)
		if err != nil {
			return err
		}
		value, err = cryptor.Decrypt(value)
		if err != nil {
			return err
		}
		p.Set(path, value)
	}
	return nil
}
