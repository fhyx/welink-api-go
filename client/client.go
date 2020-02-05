package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

	logger().Debugw("client Do", "method", method, "uri", uri)
	req, e := http.NewRequest(method, uri, body)
	if e != nil {
		log.Println(e, method, uri)
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
		log.Printf("client %s %s ERR %s", req.Method, req.RequestURI, e)
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		log.Printf("http code error %d, %s", resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Printf("read body ERR %s", e)
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

func (c *Client) GetJSON(uri string, obj interface{}) error {
	body, err := c.Get(uri)
	if err != nil {
		return err
	}
	err = parseResult(body, obj)
	if err != nil {
		log.Printf("GetJSON(uri %s) ERR %s", uri, err)
	}
	return err
}

func (c *Client) PostJSON(uri string, data []byte, obj interface{}) error {
	c.ctype = "application/json"
	body, err := c.Post(uri, data)
	if err != nil {
		return err
	}
	err = parseResult(body, obj)
	if err != nil {
		log.Printf("PostJSON(uri %s, %d bytes) ERR %s", uri, len(data), err)
	}
	return err
}

func parseResult(resp []byte, obj interface{}) error {
	// log.Printf("parse result: %s", string(resp))
	exErr := &Error{}
	if e := json.Unmarshal(resp, exErr); e != nil {
		log.Printf("unmarshal api err %s", e)
		return e
	}

	if exErr.Code != 0 || (exErr.Message != "SUCCESS" && exErr.Message != "OK") {
		log.Printf("apiError %s", exErr)
		return exErr
	}

	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return e
	}

	return nil
}
