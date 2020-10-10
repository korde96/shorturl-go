package store

import (
	"crypto/md5"
	"errors"
	"fmt"
	"shorturl-go/config"

	as "github.com/aerospike/aerospike-client-go"
)

type Store interface {
	PutSlug(string, string) error
	GetSlug(string) (string, error)
	GetIfExists(v string) (string, error)
}

// type badgerstore struct {
// 	db *badger.DB
// }

// func NewBadgerStore(db *badger.DB) Store {
// 	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &badgerstore{db}
// }

// func (m *badgerstore) Put(k, v string) error {

// 	_, err := m.Get(k)
// 	if err == nil {
// 		return errors.New("Already Exists")
// 	}

// 	err = m.db.Update(func(txn *badger.Txn) error {
// 		e := badger.NewEntry([]byte(k), []byte(v)).WithTTL(5 * time.Minute)
// 		err := txn.SetEntry(e)
// 		return err
// 	})
// 	return err
// }

// func (m *badgerstore) Get(k string) (string, error) {
// 	var value []byte
// 	err := m.db.View(func(txn *badger.Txn) error {

// 		v, err := txn.Get([]byte(k))
// 		if err != nil {
// 			if err == badger.ErrKeyNotFound {
// 				return errors.New("Not Found")
// 			}
// 			return err
// 		}
// 		value, err = v.ValueCopy([]byte{})
// 		return err
// 	})
// 	return string(value), err
// }

var (
	slugSet string = "slug"
	urlHash string = "hash"
)

var (
	ASError     = errors.New("Aerospike Error")
	ErrNotFound = errors.New("Not Found")
)

type aerospike struct {
	as  *as.Client
	cfg config.ASConfig
}

func NewASStore(as *as.Client, cfg config.ASConfig) Store {
	return &aerospike{as, cfg}
}

func (a *aerospike) PutSlug(k, v string) error {
	fmt.Printf("Putting %s,%s\n", k, v)
	key, err := as.NewKey(a.cfg.Namespace, slugSet, k)
	if err != nil {
		return ASError
	}
	policy := as.NewWritePolicy(0, uint32(a.cfg.TTL))
	policy.RecordExistsAction = as.CREATE_ONLY

	err = a.as.PutObject(policy, key, &Slug{k, v})
	if err != nil {
		fmt.Printf("Could not put %s,%s\n", k, v)
		return ASError
	}

	kh, err := as.NewKey(a.cfg.Namespace, urlHash, md5.Sum([]byte(v)))
	if err != nil {
		return ASError
	}
	err = a.as.PutObject(as.NewWritePolicy(0, uint32(a.cfg.TTL)), kh, &Slug{k, v})
	if err != nil {
		return ASError
	}
	return nil
}

func (a *aerospike) GetSlug(k string) (string, error) {
	key, err := as.NewKey(a.cfg.Namespace, slugSet, k)
	if err != nil {
		return "", ASError
	}
	record := &Slug{}
	err = a.as.GetObject(as.NewPolicy(), key, record)
	if err != nil {
		return "", ErrNotFound
	}
	return record.Url, err
}

func (a *aerospike) GetIfExists(v string) (string, error) {
	key, err := as.NewKey(a.cfg.Namespace, urlHash, md5.Sum([]byte(v)))
	if err != nil {
		return "", ASError
	}
	record := &Slug{}
	err = a.as.GetObject(as.NewPolicy(), key, record)
	return record.Slug, err
}

type Slug struct {
	Slug string `as:"slug"`
	Url  string `as:"url"`
}
