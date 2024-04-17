package paperless

import (
	"net/http"
	"strings"
)

type Client struct {
	URL        string
	HttpClient *http.Client

	username string
	token    string
}

// NewClient creates a new PaperlessClient using the given URL and credentials.
// If using token auth, `username` parameter can be left empty.
func NewClient(url, username, passwordOrToken string) *Client {
	return &Client{
		URL:        strings.TrimSuffix(url, "/"),
		HttpClient: http.DefaultClient,
		username:   username,
		token:      passwordOrToken,
	}
}

func (clt *Client) setAuth(req *http.Request) {
	if clt.username == "" {
		req.Header.Set("Authorization", "Token "+clt.token)
	} else {
		req.SetBasicAuth(clt.username, clt.token)
	}
}
