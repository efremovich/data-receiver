package httputil

import (
	"encoding/base64"
	"time"

	"github.com/valyala/fasthttp"
)

func SendHTTPRequest(method string, url string, body []byte, headers map[string]string, login, pass string, timeout time.Duration) (int, []byte, error) {

	req := fasthttp.AcquireRequest()

	if login != "" && pass != "" {
		token := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
		req.Header.Add("Authorization", "Basic "+token)
	}

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentLength(len(body))

	req.SetBody(body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp := fasthttp.AcquireResponse()

	err := fasthttp.DoTimeout(req, resp, timeout)
	if err != nil {
		return 0, []byte{}, err
	}

	respBody := resp.Body()
	statusCode := resp.StatusCode()

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return statusCode, respBody, nil
}
