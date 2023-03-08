package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	tendpoint "github.com/RangelReale/go-kit-typed/endpoint"
)

func makeUppercaseEndpoint(svc StringService) tendpoint.Endpoint[uppercaseRequest, uppercaseResponse] {
	return func(ctx context.Context, req uppercaseRequest) (uppercaseResponse, error) {
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) tendpoint.Endpoint[countRequest, countResponse] {
	return func(ctx context.Context, req countRequest) (countResponse, error) {
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (uppercaseRequest, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return uppercaseRequest{}, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (countRequest, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return countRequest{}, err
	}
	return request, nil
}

func decodeUppercaseResponse(_ context.Context, r *http.Response) (uppercaseResponse, error) {
	var response uppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return uppercaseResponse{}, err
	}
	return response, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}
