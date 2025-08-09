package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPError struct {
	Code int
	Body string
}

func (e *HTTPError) Error() string { return fmt.Sprintf("http error: %d %s", e.Code, e.Body) }
func IsHTTPStatus(err error, code int) bool {
	if he, ok := err.(*HTTPError); ok {
		return he.Code == code
	}
	return false
}

func DoJSON(ctx context.Context, cli *HTTP, req *http.Request, out any) error {
	resp, err := cli.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{Code: resp.StatusCode, Body: string(preview(b))}
	}
	if out == nil || len(b) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode json into %T: %w; body: %s", out, err, preview(b))
	}
	return nil
}

func preview(b []byte) string {
	if len(b) > 2048 {
		b = b[:2048]
	}
	// جلو HTML:
	s := strings.TrimSpace(string(b))
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
