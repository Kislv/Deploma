package telegram

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage/files"
)

type Client struct {
	host     string
	basePath string
	fileBasePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	getFileMethod  = "getFile"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		fileBasePath: newFileBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func newFileBasePath(token string) string {
	return path.Join("file", "bot" + token)
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) File (ctx context.Context, fileId string) (filePath string, err error) {
	defer func() { err = e.WrapIfErr("can't get file", err) }()

	q := url.Values{}
	q.Add("file_id", fileId)

	data, err := c.doRequest(ctx, getFileMethod, q)
	if err != nil {
		return  "", err
	}

	var res FileResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return  "", err
	}

	if err := files.IsValidExtension(res.Result.FilePath); err  != nil{
		return "", err
	}

	return  res.Result.FilePath, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	
	req.URL.RawQuery = query.Encode()
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) DownloadFile (ctx context.Context, filePath string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't download file", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.fileBasePath, filePath),
	}
	println("DownloadFile URL:" + u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

