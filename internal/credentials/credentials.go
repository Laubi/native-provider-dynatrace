package credentials

import (
	"encoding/json"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/pkg/errors"
)

const (
	errCredentials = "cannot unmarshal credentials"
)

type Credentials struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func Unmarshal(credentialBytes []byte) (*settings.Credentials, error) {
	c := Credentials{}
	if err := json.Unmarshal(credentialBytes, &c); err != nil {
		return nil, errors.Wrap(err, errCredentials)
	}

	return &settings.Credentials{
		URL:   c.Url,
		Token: c.Token,
	}, nil
}
