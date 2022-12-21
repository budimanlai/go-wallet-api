package library

import (
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

func SendNotif(url_ string, body url.Values) (*fasthttp.Response, error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url_)
	req.Header.DisableNormalizing()
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	req.SetBodyString(body.Encode())

	respClone := &fasthttp.Response{}
	e := fasthttp.DoTimeout(req, resp, 60*time.Second)
	resp.CopyTo(respClone)

	return respClone, e
}
