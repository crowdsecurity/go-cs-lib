package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/crowdsecurity/go-cs-lib/downloader"
)

func main() {
	url := "https://raw.githubusercontent.com/crowdsecurity/crowdsec/master/cmd/crowdsec-cli/alerts.go"
	myfile := "/tmp/foo/bar/baz/myfile"

	client := http.Client{}

	log.SetLevel(log.DebugLevel)

	d := downloader.New().
		WithLogger(log.StandardLogger().WithFields(log.Fields{"url": url})).
		WithHTTPClient(&client).
		ToFile(myfile).
		WithMakeDirs(true).
		//		WithETagFile(myfile+".etag").
		WithETagFn(downloader.SHA256).
		//		IfModifiedSince().
		//		WithLastModified().
		WithShelfLife(7*24*time.Hour).
		WithMode(0o640).
		LimitDownloadSize(1024*100).
		CompareContent().
		VerifyHash("sha256", "6ed6e688e3e4c916ec310600e10d16883ee7a03c0e4c46e227ae5459902bf029")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//	go func() {
	//		time.Sleep(5 * time.Second)
	//		cancel()
	//	}()

	downloaded, err := d.Download(ctx, url)
	if err != nil {
		fmt.Print(err, "\n")
	}

	fmt.Println("Downloaded: ", downloaded)
}
