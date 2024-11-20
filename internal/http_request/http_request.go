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

func doRequest(url string, httpVerb string, headers []map[string]string) *HttpRequest {
	req, err := http.NewRequest(httpVerb, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, header := range headers {
		req.Header.Add(header["key"], header["value"])
	}

	defer resp.Body.Close()
	return &HttpRequest{Body: string(body), Resp: resp}
}

func FetchUrl(url string, httpVerb string, headers []map[string]string) *HttpRequest {

	var req *HttpRequest

	switch httpVerb {
	case "GET":
		req = doRequest(url, "GET", headers)
	case "POST":
		req = doRequest(url, "POST", headers)
	case "PUT":
		req = doRequest(url, "PUT", headers)
	case "DELETE":
		req = doRequest(url, "DELETE", headers)
	default:
		log.Fatal("Invalid or not found HTTP Verb")
		return nil
	}

	return &HttpRequest{Body: string(req.Body), Resp: req.Resp}
}
