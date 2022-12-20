package slauth

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		StatusCode: 0,
	}
}

func (r *HttpResponse) SetStatusCode(code int) *HttpResponse {
	r.StatusCode = code
	return r
}

func (r *HttpResponse) SetBody(body []byte) *HttpResponse {
	r.Body = body
	return r
}
