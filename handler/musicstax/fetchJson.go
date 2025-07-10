package musicstax

import (
	"fmt"
	"io"
	"net/http"
)

func FetchJSON(id, cf, ua string) ([]byte, error) {
	url := fmt.Sprintf("https://musicstax.com/track/%s.json?similar=true", id)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", ua)
	req.AddCookie(&http.Cookie{
		Name:  "cf_clearance",
		Value: cf,
	})
	// cli := &http.Client{Timeout: 20 * time.Second}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return io.ReadAll(res.Body)
}