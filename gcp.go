package nvault

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	kms "google.golang.org/api/cloudkms/v1"
)

type GcpCryptor struct {
	GcpConfig
}

type GcpConfig struct {
	GcpKmsResourceID  string
	GcpCredentialFile string
}

func (c *GcpCryptor) Encrypt(value interface{}) (interface{}, error) {
	if c.GcpKmsResourceID == "" {
		return nil, errors.New("missing Gcp KMS Resource ID")
	}
	strvalue := fmt.Sprintf("%v", value)

	svc, err := serviceGcp(&c.GcpConfig)
	if err != nil {
		return nil, err
	}

	response, err := svc.Projects.Locations.KeyRings.CryptoKeys.Encrypt(c.GcpKmsResourceID, &kms.EncryptRequest{
		Plaintext: strvalue,
	}).Do()
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(response.Ciphertext))
	return encoded, nil
}

func (c *GcpCryptor) Decrypt(value interface{}) (interface{}, error) {
	strvalue := fmt.Sprintf("%v", value)

	decoded, err := base64.StdEncoding.DecodeString(strvalue)
	if err != nil {
		return nil, err
	}

	svc, err := serviceGcp(&c.GcpConfig)

	response, err := svc.Projects.Locations.KeyRings.CryptoKeys.Decrypt(c.GcpKmsResourceID, &kms.DecryptRequest{
		Ciphertext: string(decoded),
	}).Do()
	if err != nil {
		return value, nil
	}

	return string(response.Plaintext), nil
}

func serviceGcp(c *GcpConfig) (*kms.Service, error) {
	client, err := createGcpClient(c)
	if err != nil {
		return nil, err
	}
	return kms.New(client)
}

func createGcpClient(c *GcpConfig) (client *http.Client, err error) {
	ctx := context.Background()
	if c.GcpCredentialFile != "" {
		data, err := ioutil.ReadFile(c.GcpCredentialFile)
		if err != nil {
			return nil, err
		}
		creds, err := google.CredentialsFromJSON(ctx, data, kms.CloudPlatformScope)
		if err != nil {
			return nil, err
		}
		client = oauth2.NewClient(ctx, creds.TokenSource)
		return client, nil
	} else {
		return google.DefaultClient(ctx, kms.CloudPlatformScope)
	}
}
