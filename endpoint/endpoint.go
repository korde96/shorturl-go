package endpoint

import (
	"context"
	"shorturl-go/svc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	AddURLEndpoint endpoint.Endpoint
	GetURLEndpoint endpoint.Endpoint
}

func MakeShortURLEndpoints(svc svc.URLSvc, logger log.Logger) Endpoints {
	addURLEndpoint := makeAddURLEndpoint(svc)
	addURLEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "addURL"))(addURLEndpoint)

	getURLEndpoint := makeGetURLEndpoint(svc)
	getURLEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "getURL"))(getURLEndpoint)

	return Endpoints{
		AddURLEndpoint: addURLEndpoint,
		GetURLEndpoint: getURLEndpoint,
	}
}

func makeAddURLEndpoint(svc svc.URLSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddURLRequest)
		surl, err := svc.ShortenURL(req.URL)
		return AddURLResponse{surl}, err
	}
}

func makeGetURLEndpoint(svc svc.URLSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetURLRequest)
		url, err := svc.GetURL(req.SURL)
		return GetURLResponse{url}, err
	}
}

type AddURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type AddURLResponse struct {
	SURL string `json:"url"`
}

type GetURLRequest struct {
	SURL string `json:"url"`
}

type GetURLResponse struct {
	URL string `json:"url"`
}
