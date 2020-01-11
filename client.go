package melian

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type Client struct {
	DSN         string
	Dist        string
	Environment string
	Release     string
	ServerName  string
	HTTPClient  *http.Client
}

func (c *Client) Capture(ctx context.Context, err error) error {
	e := generateEvent(err)
	populateClientData(c, &e)
	applyContext(ctx, &e)

	url, cerr := dsnToURL(c.DSN)
	if cerr != nil {
		return cerr
	}

	body, cerr := json.Marshal(e)
	if cerr != nil {
		return cerr
	}

	req, cerr := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if cerr != nil {
		return cerr
	}

	cerr = addAuthHeader(c.DSN, req)
	if cerr != nil {
		return cerr
	}

	res, cerr := c.HTTPClient.Do(req)
	if cerr != nil {
		return cerr
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Got error from sentry: %s", res.Status)
	}

	return nil
}

func addAuthHeader(dsn string, req *http.Request) error {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return err
	}

	if parsedURL.User == nil {
		return fmt.Errorf("invalid DSN")
	}

	// PublicKey
	publicKey := parsedURL.User.Username()
	if publicKey == "" {
		return fmt.Errorf("invalid DSN")
	}

	// SecretKey
	secretKey, _ := parsedURL.User.Password()

	auth := fmt.Sprintf(
		"Sentry sentry_version=%s, sentry_timestamp=%d, sentry_client=sentry.go/%s, sentry_key=%s",
		"7",
		time.Now().Unix(),
		"0.4.0", // we're totally the go sdk >.>
		publicKey,
	)

	if secretKey != "" {
		auth = fmt.Sprintf("%s, sentry_secret=%s", auth, secretKey)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Sentry-Auth", auth)
	return nil
}

func dsnToURL(dsn string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	if len(parsedURL.Path) == 0 {
		return "", fmt.Errorf("invalid DSN")
	}

	projectID := parsedURL.Path[1:]

	parsedURL.Path = fmt.Sprintf("/api/%s/store/", projectID)
	parsedURL.User = nil
	return parsedURL.String(), nil
}

func generateEvent(err error) event {
	return event{
		EventID:   uuid(),
		Platform:  "go",
		Logger:    "melian",
		Level:     "error",
		Message:   err.Error(),
		Timestamp: time.Now().Unix(),
		Exception: []exception{
			exception{
				Type:       reflect.TypeOf(err).String(),
				Value:      err.Error(),
				Stacktrace: newStacktrace(),
			},
		},
	}
}

func uuid() string {
	id := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, id)
	id[6] &= 0x0F // clear version
	id[6] |= 0x40 // set version to 4 (random uuid)
	id[8] &= 0x3F // clear variant
	id[8] |= 0x80 // set to IETF variant
	return hex.EncodeToString(id)
}

func populateClientData(c *Client, e *event) {
	e.Dist = c.Dist
	e.Environment = c.Environment
	e.Release = c.Release
	e.ServerName = c.ServerName
}

func applyContext(ctx context.Context, e *event) {
	if extra, ok := getExtra(ctx); ok {
		e.Extra = extra
	}

	if tags, ok := getTags(ctx); ok {
		e.Tags = tags
	}

	if request, ok := getRequest(ctx); ok {
		e.Request = request
	}

	if transaction, ok := getTransaction(ctx); ok {
		e.Transaction = transaction
	}
}
