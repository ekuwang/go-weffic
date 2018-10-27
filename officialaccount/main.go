package officialaccount

import "fmt"

type OfficialAccount struct {
	config *Config
	server *server
}

func New(config *Config) (*OfficialAccount, error) {
	if config.AppID == "" {
		return nil, fmt.Errorf("need appid")
	}

	if config.AppSecret == "" {
		return nil, fmt.Errorf("need appsecret")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("need token")
	}

	return &OfficialAccount{
		config: config,
	}, nil
}
