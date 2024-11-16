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

func DoRequest(url string, httpVerb string) *HttpRequest {
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

	defer resp.Body.Close()
	return &HttpRequest{Body: string(body), Resp: resp}
}

func FetchUrl(url string, httpVerb string) *HttpRequest {

	var req *HttpRequest

	switch httpVerb {
	case "GET":
		req = DoRequest(url, "GET")
	case "POST":
		req = DoRequest(url, "POST")
	case "PUT":
		req = DoRequest(url, "PUT")
	case "DELETE":
		req = DoRequest(url, "DELETE")
	default:
		log.Fatal("Invalid or not found HTTP Verb")
		return nil
	}

	return &HttpRequest{Body: string(req.Body), Resp: req.Resp}
}
