package tbapi

import (
	"errors"
	"net/url"
)

type Builder struct {
	hostname string
	username string
	password string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithHostname(hostname string) *Builder {
	b.hostname = hostname
	return b
}

func (b *Builder) WithUsername(username string) *Builder {
	b.username = username
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.password = password
	return b
}

func (b *Builder) Build() (*TabroomApi, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}
	parsed, err := url.Parse(b.hostname)
	if err != nil {
		return nil, err
	}
	return &TabroomApi{
		username: b.username,
		password: b.password,
		client:   newDefaultHttpClient(*parsed),
	}, nil
}

func (b *Builder) validate() error {
	if b.hostname == "" {
		return errors.New("missing API URL in builder")
	}
	if b.username == "" {
		return errors.New("missing username in builder")
	}
	if b.password == "" {
		return errors.New("missing password in builder")
	}
	return nil
}
