package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	authHeader = "x-wlk-Authorization"
)

type Client struct {
	*TokenHolder
	httpClient *http.Client
	ctype      string
}

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// DisableCompression: true,
	}
)

func NewClient(urlToken string) *Client {
	hc := &http.Client{Transport: tr}
	return &Client{
		httpClient:  hc,
		TokenHolder: NewTokenHolder(urlToken),
	}
}

func (c *Client) SetContentType(ctype string) {
	if ctype != "" {
		c.ctype = ctype
	}
}

func (c *Client) Do(method, uri string, body io.Reader) ([]byte, error) {
	token, err := c.GetAuthToken()
	if err != nil {
		return nil, err
	}

	logger().Debugw("client Do", method, uri)
	req, e := http.NewRequest(method, uri, body)
	if e != nil {
		logger().Infow("client do fail", "err", e)
		return nil, e
	}

	if method == "POST" {
		if c.ctype != "" {
			req.Header.Set("Content-Type", c.ctype)
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	req.Header.Set(authHeader, token)

	return doRequest(c.httpClient, req)
}

func doRequest(hc *http.Client, req *http.Request) ([]byte, error) {
	resp, e := hc.Do(req)
	if e != nil {
		logger().Infow("do request fail", req.Method, req.RequestURI, "err", e)
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		logger().Infow("http fail", "code", resp.StatusCode, "status", resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logger().Infow("read fail", "err", e)
		return nil, e
	}
	logger().Debugw("resp", "body", string(rbody))

	return rbody, nil

}

func (c *Client) Get(uri string) ([]byte, error) {
	return c.Do("GET", uri, nil)
}

func (c *Client) Post(uri string, data []byte) ([]byte, error) {
	return c.Do("POST", uri, bytes.NewReader(data))
}

// GetJSON ...
func (c *Client) GetJSON(uri string, obj interface{}) error {
	body, err := c.Get(uri)
	if err != nil {
		logger().Infow("get json fail", "uri", uri, "err", err)
		return err
	}
	err = parseResult(body, obj)
	if err != nil {
		logger().Infow("parse resp fail", "uri", uri, "body", string(body), "err", err)
	}
	return err
}

// PostJSON ...
func (c *Client) PostJSON(uri string, data []byte, obj interface{}) error {
	logger().Debugw("post json", "data", string(data))

	c.ctype = "application/json"
	body, err := c.Post(uri, data)
	if err != nil {
		logger().Infow("post fail", "uri", uri, "data", len(data), "err", err)
		return err
	}

	err = parseResult(body, obj)
	if err != nil {
		logger().Infow("parse resp fail", "uri", uri, "data", len(data), "err", err)
	}
	return err
}

func parseResult(resp []byte, obj interface{}) error {
	exErr := &Error{}
	if e := json.Unmarshal(resp, exErr); e != nil {
		logger().Infow("json unmarshal fail", "err", e)
		return e
	}

	if exErr.Code != 0 || (exErr.Message != "SUCCESS" && exErr.Message != "OK" && exErr.Message != "ok") {
		logger().Infow("parse result fail", "err", exErr)
		return exErr
	}

	if e := json.Unmarshal(resp, obj); e != nil {
		logger().Infow("json unmarshal fail", "err", e)
		return e
	}

	return nil
}
