package wx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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
		log.Fatal("http.NewRequest error: ", err)
		return err
	}
	c.req = req

	err = c.prepareQueries()
	if err != nil {
		log.Print("prepare http request failed: ", err)
		return err
	}

	err = c.request(value)
	if err != nil {
		log.Print("http request failed: ", err)
		return err
	}

	return nil
}

func (c *HttpClient) DoPost(body io.Reader, value interface{}) (err error) {
	req, err := http.NewRequest("POST", c.Path, body)
	if err != nil {
		log.Fatal("http.NewRequest error: ", err)
		return err
	}
	c.req = req
	req.Header.Set("content_type", c.ContentType)

	err = c.prepareQueries()
	if err != nil {
		log.Print("prepare http request failed: ", err)
		return err
	}

	err = c.request(value)
	if err != nil {
		log.Print("http request failed: ", err)
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
	log.Print("req: ", c.req)

	resp, err := client.Do(c.req)
	if err != nil {
		log.Print("http request error: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = HttpError{
			State: resp.StatusCode,
		}
		log.Print("http request not ok: ", err)
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("read response body error: ", err)
		return err
	}

	err = handleRespError(data)
	if err != nil {
		log.Print("http response error: ", err)
		return err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		log.Print("unmarshal response body error: ", err)
		return err
	}

	return nil
}
