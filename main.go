package main

import (
	"./core"
	"flag"
)

func main() {
	var filePath string
	var nick string
	flag.StringVar(&filePath, "PATH", "//", "abc/d/e")
	flag.StringVar(&nick, "NICK", "###", "ABCDE")
	flag.Parse()
	if filePath != "//" {
		core.Send(filePath) //To send the file needs filePath
	}
	if nick != "###" {
		core.Recieve(nick) //To recieve the file needs Nick
	}
}
