package client

import (
	"bytes"
	"encoding/json"
	"errors"
	// "fmt"
	"log"
	"net/http"
	"time"
)

type tokenReq struct {
	CorpID     string `json:"client_id"`
	CorpSecret string `json:"client_secret"`
}

// Token ...
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Error
}

// TokenHolder ...
type TokenHolder struct {
	currToken  *Token
	uri        string
	method     string
	apiAuths   string
	corpID     string
	corpSecret string
	expiresAt  int64
}

var (
	errEmptyAuths = errors.New("empty auth string or corpID and corpSecret")
)

func NewTokenHolder(uri string) *TokenHolder {
	return &TokenHolder{
		uri:    uri,
		method: "POST",
	}
}

func (th *TokenHolder) SetAuth(auths string) {
	th.apiAuths = auths
}

func (th *TokenHolder) SetCorp(id, secret string) {
	th.corpID = id
	th.corpSecret = secret
}

func (th *TokenHolder) Expired() bool {
	return th.expiresAt < time.Now().Unix()
}

func (th *TokenHolder) Valid() bool {
	if th.currToken == nil {
		return false
	}
	return !th.Expired()
}

func (th *TokenHolder) GetAuthToken() (token string, err error) {
	if !th.Valid() {
		logger().Debugw("token is nil or expired, refreshing it")
		th.currToken, err = th.requestToken()
		if err != nil {
			return "", err
		}
		// log.Print("got token", th.currToken)
		th.expiresAt = time.Now().Unix() + th.currToken.ExpiresIn
	}
	token = th.currToken.AccessToken
	return
}

func (th *TokenHolder) requestToken() (token *Token, err error) {
	body, _ := json.Marshal(&tokenReq{th.corpID, th.corpSecret})
	req, err := http.NewRequest("POST", th.uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	hc := &http.Client{Transport: tr}
	var resp []byte
	resp, err = doRequest(hc, req)

	if err != nil {
		log.Printf(" err %s", err)
		return
	}

	obj := &Token{}
	err = json.Unmarshal(resp, obj)
	if err != nil {
		log.Printf("unmarshal err %s", err)
		return
	}
	token = obj

	return
}
