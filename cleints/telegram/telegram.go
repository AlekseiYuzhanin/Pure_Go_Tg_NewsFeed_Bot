package telegram

import (
	err2 "awesomeProject4/lib/err"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Set("chat_id", strconv.Itoa(chatID))
	q.Set("text", text)

	_, err := c.Request(sendMessageMethod, q)
	if err != nil {
		return err2.Wrap(err, "failed to send message")
	}
	return nil
}

func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	defer func() { err = err2.WrapIfErr(err, "cannot get updates") }()
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.Request(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) Request(method string, query url.Values) (data []byte, err error) {
	defer func() { err = err2.WrapIfErr(err, "cant do request") }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
