//
// client.go
// Copyright (C) 2018 toraxie <toraxie@tencent.com>
//
// Distributed under terms of the Tencent license.
//

package tgdb

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	bolt "github.com/mindstand/go-bolt"
	"github.com/mindstand/go-bolt/bolt_mode"
)

var (
	public *rsa.PublicKey
)

const (
	tgdbPubKey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCXgctF+noMbOXv5c8hVef4vE6eK/mwnzEp
MoHpZJjphuyTcowsD9DBWe+PrqZOEm1PomzD/TxOn9eMhn9O+w3yuuv8fipMs6OjU6Y+rKLr
GJ8aCZpzd7Y5eewQS/0hOGmWtQlLdayJaUT0B0Fpz3yhR7u7vtVhKwbCYYfHxhg1PQIDAQAB
-----END PUBLIC KEY-----
`

	ReadMode  = bolt_mode.ReadMode
	WriteMode = bolt_mode.WriteMode
)

type Client struct {
	pass string // encrypt pass

	c bolt.IClient
	d bolt.IDriver
}

type AccessMode = bolt_mode.AccessMode

func initPubKey() {
	block, _ := pem.Decode([]byte(tgdbPubKey))
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("init public key failed %w", err))
	}
	public = key.(*rsa.PublicKey)
}

func init() {
	initPubKey()
}

func encryptPassword(pass string) (string, error) {
	enc, err := rsa.EncryptPKCS1v15(rand.Reader, public, []byte(pass))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(enc), nil

}

func New(host string, port int, user, pass string) (*Client, error) {
	ret := &Client{}
	var err error
	ret.pass, err = encryptPassword(pass)
	if err != nil {
		return nil, err
	}
	ret.c, err = bolt.NewClient(bolt.WithBasicAuth(user, ret.pass), bolt.WithHostPort(host, port))
	if err != nil {
		return nil, err
	}
	ret.d, err = ret.c.NewDriver()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) Open(db string, mode AccessMode) (*Session, error) {
	ret := &Session{}
	var err error
	ret.conn, err = c.d.Open(mode)
	if err != nil {
		return nil, err
	}
	ret.db = db
	ret.p = c
	return ret, nil
}
