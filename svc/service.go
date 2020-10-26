package svc

import (
	"errors"
	"shorturl-go/internal/store"

	"github.com/go-kit/kit/log"
)

var ErrSlugGenFailed = errors.New("Failed to Generate Short URL")

type URLSvc interface {
	ShortenURL(url string, slugProvider func() (string, error)) (string, error)
	GetURL(url string) (string, error)
}

type urlsvc struct {
	store  store.Store
	logger log.Logger
}

func NewURLSvc(store store.Store, logger log.Logger) URLSvc {
	logger = log.With(logger, "service", "urlSvc")
	return &urlsvc{store, logger}
}

func (svc *urlsvc) ShortenURL(url string, slugProvider func() (string, error)) (string, error) {
	slug, err := svc.store.GetIfExists(url)
	if err == nil {
		svc.logger.Log("msg", "URL already exists in store")
		return slug, nil
	}

	svc.logger.Log("msg", "Generating Short URL")
	var surl string
	if err = retryGen(func() error {
		if surl, err = slugProvider(); err == nil {
			if storeError := svc.store.PutSlug(surl, url); storeError == nil {
				return nil
			} else {
				return storeError
			}
		} else {
			return err
		}
	}, 3); err == nil {
		return surl, nil
	}

	return "", ErrSlugGenFailed
}

func (svc *urlsvc) GetURL(surl string) (string, error) {
	return svc.store.GetSlug(surl)
}

func retryGen(f func() error, n int) error {
	var err error
	for n > 0 {
		if err = f(); err != nil {
			n--
		} else {
			return nil
		}
	}
	return err
}
