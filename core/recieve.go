package core

import (
	"bufio"
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func findServer(name string, entries chan *zeroconf.ServiceEntry) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize")
	}
	// Context for 10 seconds , if after 10 seconds i can't find the server i will stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = resolver.Lookup(ctx, name, "_share._http._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Cant find a server")
	}
	<-ctx.Done()
	return
}

// Recieve the file
func Recieve(name string) {
	var datalen int
	var entry *zeroconf.ServiceEntry
	entries := make(chan *zeroconf.ServiceEntry)
	go findServer(GetHash(name), entries)
	name = url.PathEscape(name)
	entry = <-entries // blocks till it finds the server
	if entry == nil {
		log.Fatalln("can't find server in 10 seconds check if the server is on")
	}
	addr := entry.AddrIPv4[0]
	port := entry.Port

	urlPath := fmt.Sprintf("http://%v.%v.%v.%v:%d/%s", addr[0], addr[1], addr[2], addr[3], port, name)
	resp, err := http.Get(urlPath)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != 200 {
		log.Fatalln("File not correct status code:", resp.StatusCode)
	} else {
		output := resp.Header.Get("Name")
		output = SafeFilename(output) //safefilename appends a number to the name of file if it already exists
		of, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		fmt.Sscanf(resp.Header.Get("len"), "%d", &datalen)
		respBuf := bufio.NewReader(resp.Body)
		ofBuf := bufio.NewWriter(of)
		bar := pb.New(datalen).SetUnits(pb.U_BYTES)
		bar.Start()
		barWriter := io.MultiWriter(ofBuf, bar)
		io.Copy(barWriter, respBuf)
		bar.Finish()
	}
}
