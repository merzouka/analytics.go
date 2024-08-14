package main

import "time"

type Response struct {
	Data      interface{}   `json:"data"`
	Errors    []string      `json:"errors"`
	Source    string        `json:"source"`
	Succeeded bool          `json:"succeeded"`
	Duration  time.Duration `json:"duration"`
}

func New(source string) Response {
	return Response{
		Source: source,
	}
}

func (r Response) SetData(data func() (interface{}, error)) Response {
	var err error

	start := time.Now()
	r.Data, err = data()
	r.Succeeded = err == nil

	if err != nil {
		r.Errors = append(r.Errors, err.Error())
	}
	r.Duration = time.Now().Sub(start)

	return r
}

func (r Response) AddErrors(errors ...error) Response {
	for _, err := range errors {
		r.Errors = append(r.Errors, err.Error())
	}

	return r
}

type SSEResponse struct {
	Done     bool          `json:"done"`
	Source   string        `json:"source"`
	Data     interface{}   `json:"data"`
	Duration string        `json:"duration"`
    Success  bool          `json:"success"`
}

func NewSSE(source string) SSEResponse {
	return SSEResponse{
		Done:   false,
		Source: source,
	}
}

func (r SSEResponse) SetDuration(t time.Duration) SSEResponse {
	r.Duration = t.String()
	return r
}

func (r SSEResponse) SetData(data interface{}, success bool) SSEResponse {
	r.Data = data
    r.Success = success
	return r
}

func (r SSEResponse) Final() SSEResponse {
	r.Done = true
	return r
}
