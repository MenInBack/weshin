package wx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultTimeout = 10
)

type HttpClient struct {
	Path        string
	Parameters  []QueryParameter
	Timeout     int
	ContentType string
	req         *http.Request
}

type QueryParameter struct {
	Key   string
	Value string
}

// to reserve the order of parameters
func (c HttpClient) Get(value interface{}) error {
	req, err := http.NewRequest("GET", c.Path, nil)
	if err != nil {
		return err
	}
	c.req = req

	err = c.prepareQueries()
	if err != nil {
		return err
	}

	err = c.request(value)
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpClient) DoPost(body io.Reader, value interface{}) (err error) {
	req, err := http.NewRequest("POST", c.Path, body)
	if err != nil {
		return err
	}
	c.req = req
	req.Header.Set("content_type", c.ContentType)

	err = c.prepareQueries()
	if err != nil {
		return err
	}

	err = c.request(value)
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpClient) prepareQueries() error {
	var q bytes.Buffer
	for i, p := range c.Parameters {
		if i > 0 {
			q.WriteString("&")
		}
		q.WriteString(p.Key)
		q.WriteString("=")
		q.WriteString(p.Value)
	}

	c.req.URL.RawQuery = q.String()
	return nil
}

func (c *HttpClient) request(value interface{}) error {
	client := http.Client{
		Timeout: func() time.Duration {
			if c.Timeout > 0 {
				return time.Duration(c.Timeout) * time.Second
			}
			return defaultTimeout * time.Second
		}(),
	}

	resp, err := client.Do(c.req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = HttpError{
			State: resp.StatusCode,
		}
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = handleRespError(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}
