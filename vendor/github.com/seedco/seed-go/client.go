package seed

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ApiBase = "https://api.seed.co/v1/public"
)

type Client struct {
	accessToken   string
	clientVersion string
	httpClient    *http.Client
}

// New creates a new Seed client given an access token
func New(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
	}
}

// SetClientVersion sets the client version that each request will be made with
func (c *Client) SetClientVersion(clientVersion string) {
	c.clientVersion = clientVersion
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	if c.clientVersion != "" {
		req.Header.Set("Client-Version-Id", c.clientVersion)
	}
	return c.httpClient.Do(req)
}

type Pages struct {
	Next     PaginationParams `json:"next"`
	Previous PaginationParams `json:"previous"`
}

type PaginationParams struct {
	Offset int `json:"offest"`
	Limit  int `json:"limit"`
}

func (p PaginationParams) MarshalJSON() ([]byte, error) {
	return []byte(p.Encode()), nil
}

func (p *PaginationParams) UnmarshalJSON(d []byte) error {
	var err error
	s := string(d)
	params := strings.Split(s, "&")

	for _, param := range params {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			continue
		}
		switch split[0] {
		case "limit":
			var limit int
			if limit, err = strconv.Atoi(split[1]); err != nil {
				return err
			}
			p.Limit = limit
		case "offset":
			var offset int
			if offset, err = strconv.Atoi(split[1]); err != nil {
				return err
			}
			p.Offset = offset
		}
	}
	return nil
}

func (p PaginationParams) Encode() string {
	u := url.Values{}
	if p.Offset > 0 {
		u.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		u.Set("limit", strconv.Itoa(p.Limit))
	}

	return u.Encode()
}

type ErrorList []map[string]string

func (e ErrorList) Error() string {
	buf := bytes.NewBufferString("")
	for _, e := range e {
		errorString := fmt.Sprintf("field: %s, message: %s", e["field"], e["message"])
		buf.WriteString(errorString)
		buf.WriteString("\n")
	}
	return buf.String()
}
