package http_request

import (
	"io"
	"log"
	"net/http"
)

type HttpRequest struct {
	Body string
	Resp *http.Response
}

func doRequest(url string, httpVerb string, headers []map[string]string) (*HttpRequest, error) {
	req, err := http.NewRequest(httpVerb, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, header := range headers {
		req.Header.Add(header["key"], header["value"])
	}

	defer resp.Body.Close()
	return &HttpRequest{Body: string(body), Resp: resp}, nil
}

func FetchUrl(url string, httpVerb string, headers []map[string]string) (*HttpRequest, error) {
	req, err := doRequest(url, httpVerb, headers)
	if err != nil {
		return nil, err
	}
	return &HttpRequest{Body: string(req.Body), Resp: req.Resp}, nil
}
