package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shorturl-go/config"
	"shorturl-go/endpoint"
	"shorturl-go/internal/store"
	"shorturl-go/svc"
	"shorturl-go/transport"

	. "github.com/aerospike/aerospike-client-go"
	kitLog "github.com/go-kit/kit/log"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", config)

	logger := kitLog.NewLogfmtLogger(os.Stderr)

	// db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// store := store.NewBadgerStore(db)

	client, err := NewClient(config.ASConfig.Host, config.ASConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	store := store.NewASStore(client, config.ASConfig)

	urlService := svc.NewURLSvc(store, logger)

	endpoints := endpoint.MakeShortURLEndpoints(urlService, logger)

	handler := transport.NewHTTPHandler(endpoints)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%v", config.Server.Host, config.Server.Port), handler))
	// r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
