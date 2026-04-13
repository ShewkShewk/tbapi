package tb_api

import (
	"errors"
)

type Builder struct {
	Hostname string
	Username string
	Password string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithHostname(hostname string) *Builder {
	b.Hostname = hostname
	return b
}

func (b *Builder) WithUsername(username string) *Builder {
	b.Username = username
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.Password = password
	return b
}

func (b *Builder) Build() (error, *TabroomApi) {
	err := b.validate()
	if err != nil {
		return err, nil
	}
	return nil, &TabroomApi{}
}

func (b *Builder) validate() error {
	if b.Hostname == "" {
		return errors.New("missing API URL in builder")
	}
	if b.Username == "" {
		return errors.New("missing username in builder")
	}
	if b.Password == "" {
		return errors.New("missing password in builder")
	}
	return nil
}
