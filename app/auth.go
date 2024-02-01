package app

import (
	"encoding/json"
	"os"
	"path"

	"github.com/gomig/crypto"
	"github.com/gomig/utils"
)

type Credential struct {
	Key      string `json:"key"`
	Username string `json:"user"`
	Token    string `json:"token"`
}

type Authentications struct {
	credentials []Credential
	crypto      crypto.Crypto
	file        string
}

func (a *Authentications) Init() {
	a.credentials = make([]Credential, 0)
	a.crypto = crypto.NewCryptography("78e0cc765626542f21b0fcf71465f9cbdffe30c3855ec81692df37368b9901d6")
	if home, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		a.file = path.Join(home, ".mig")
	}
}

func (a *Authentications) Read() error {
	a.Init()
	if ok, err := utils.FileExists(a.file); err != nil {
		return err
	} else if ok {
		creds := make([]Credential, 0)
		raw, err := os.ReadFile(a.file)
		if err != nil {
			return err
		}

		decrypted, err := a.crypto.DecryptHex(string(raw))
		if err != nil {
			return err
		}

		err = json.Unmarshal(decrypted, &creds)
		if err != nil {
			return err
		}

		a.credentials = append(a.credentials, creds...)
	}
	return nil
}

func (a Authentications) Write() error {
	raw, err := json.Marshal(a.credentials)
	if err != nil {
		return err
	}

	encrypted, err := a.crypto.EncryptHEX(raw)
	if err != nil {
		return err
	}

	return os.WriteFile(a.file, []byte(encrypted), 0644)

}

func (a Authentications) Find(key string) *Credential {
	for _, cred := range a.credentials {
		if cred.Key == key {
			return &cred
		}
	}
	return nil
}

func (a *Authentications) Add(key, user, token string) {
	for x, cred := range a.credentials {
		if cred.Key == key {
			a.credentials[x].Key = key
			a.credentials[x].Username = user
			a.credentials[x].Token = token
			return
		}
	}
	a.credentials = append(a.credentials, Credential{
		Key:      key,
		Username: user,
		Token:    token,
	})
}

func (a *Authentications) Delete(key string) {
	for x, cred := range a.credentials {
		if cred.Key == key {
			a.credentials = append(a.credentials[:x], a.credentials[x+1:]...)
			return
		}
	}
}

func (a *Authentications) Credentials() []Credential {
	return a.credentials
}
