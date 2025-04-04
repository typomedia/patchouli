package structs

import (
	"crypto/rand"
	"encoding/hex"
)

type Config struct {
	General  General  `json:"general"`
	Smtp     Smtp     `json:"smtp"`
	Security Security `json:"security"`
}

type General struct {
	Company  string `json:"company"`
	Email    string `json:"email"`
	Interval int    `json:"interval"`
	Hostname string `json:"hostname"`
}

type Smtp struct {
	Sender   string `form:"sender" json:"sender"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Security struct {
	CipherKey string `json:"cipher_key"`
	Generated bool   `json:"generated"`
}

func (c *Config) GenerateCipherKey() error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	c.Security.CipherKey = hex.EncodeToString(bytes)
	c.Security.Generated = true

	return nil
}

func (c *Config) GetCipherKey() (string, error) {
	if c.Security.CipherKey != "" {
		return c.Security.CipherKey, nil
	}
	err := c.GenerateCipherKey()
	if err != nil {
		return "", err
	}
	return c.Security.CipherKey, nil
}
