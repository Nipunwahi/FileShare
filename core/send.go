package core

import (
	"bufio"
	"fmt"
	"github.com/grandcat/zeroconf"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	quit     chan bool
	filePath string
	stop     chan bool
)

func handler(rw http.ResponseWriter, rq *http.Request) {
	file, err := os.Open(filePath)
	var datalen int
	if err != nil {
		panic(err)
	} else {
		rw.Header().Set("type", "byteStream")
		stat, err := file.Stat()
		if err != nil {
			panic(err)
		} else {
			datalen = int(stat.Size())
			rw.Header().Add("len", fmt.Sprintf("%d", datalen)) //sending the size of the file
			rw.Header().Add("Name", stat.Name())               //sending the name to recreate it
		}
	}
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(rw)
	bar := pb.New(datalen).SetUnits(pb.U_BYTES)
	bar.Start()
	barReader := bar.NewProxyReader(reader)
	io.Copy(writer, barReader)
	bar.Finish()
	writer.Flush()
	quit <- true //gives the quit signal to shutdown the server -> see MakeZeroConf function
	stop <- true //gives the stop signal to shutdown the method -> see Send function
}

// MakeZeroConf makes a zeroconf server so it can be found
func MakeZeroConf(hash string, port int) {
	if server, err := zeroconf.Register(hash, "_share._http._tcp", "local.", port, nil, nil); err != nil {
		panic(err)
	} else {
		<-quit
		server.Shutdown()
	}

}

// Send the file
func Send(pathtofile string) {
	quit = make(chan bool)
	stop = make(chan bool)
	port := 44444
	filePath = pathtofile
	nick := GetNick(5)
	hash := GetHash(nick)
	fmt.Println(nick) // use this with -NICK Parameter to receive this file
	go MakeZeroConf(hash, port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	} else {
		file, err := os.Stat(filePath)
		if err != nil || !file.Mode().IsRegular() {
			log.Println("Not a valid file")
			os.Exit(1)
		}
		http.HandleFunc(fmt.Sprintf("/%s", nick), handler)
		go http.Serve(listener, nil)
	}
	<-stop
}
