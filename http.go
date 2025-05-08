package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HTTPTask struct {
	Method  string
	URL     string
	Headers map[string]string
	Data    string
}

func NewHTTPTask(content string) *HTTPTask {
	lines := strings.Split(content, "\n")

	header := strings.Split(lines[0], " ")
	headers := make(map[string]string)

	data := make([]string, 0)

	headerOver := false

	for _, line := range lines[1:] {
		if headerOver {
			data = append(data, line)
		} else {
			if line == "" {
				headerOver = true
			} else {
				entry := strings.Split(line, ": ")

				headers[entry[0]] = entry[1]
			}
		}
	}

	return &HTTPTask{
		Method:  header[0],
		URL:     "https://" + headers["Host"] + header[1],
		Headers: headers,
		Data:    strings.Join(data, "\n"),
	}
}

func (h *HTTPTask) Resolve(cfg *Config) StatusEntry {
	timeout := time.Duration(cfg.HTTPTimeout) * time.Second

	resp := _request(h.Method, h.URL, h.Data, h.Headers, timeout)

	for retry := 1; retry <= cfg.HTTPRetryCount && resp.Error != ""; retry++ {
		log.Debugf("Retrying %s in %ds (%d of %d)\n", h.URL, cfg.HTTPRetryDelay, retry, cfg.HTTPRetryCount)

		time.Sleep(time.Duration(cfg.HTTPRetryDelay) * time.Second)

		resp = _request(h.Method, h.URL, h.Data, h.Headers, timeout)
	}

	return resp
}

func _request(method, url, data string, headers map[string]string, timeout time.Duration) StatusEntry {
	start := time.Now()

	body := strings.NewReader(data)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return _error(err, _time(start))
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return _error(err, _time(start))
	}

	if resp.StatusCode != 200 {
		return _error(fmt.Errorf("status code was %d instead of 200", resp.StatusCode), _time(start))
	}

	return StatusEntry{
		Operational:  true,
		Type:         "http",
		ResponseTime: _time(start),
	}
}
