package main

import (
	"./api"
)

import (
	"log"
	"net/http"
	"os"
	"path"
)

var fs api.SyncFS

func RenamePath(local string, online string) {
	file, err := os.Open(local)

	if err != nil {
		log.Fatalln(err)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		log.Fatalln("stat", err)
		return
	}
	if stat.IsDir() {
		names, err := file.Readdirnames(0)
		if err != nil {
			file.Close()
			log.Fatalln("names", err)
			return
		}
		file.Close()
		for _, v := range names {
			RenamePath(local+"/"+v, online+"/"+v)
		}
	} else {
		timestamp := stat.ModTime().Format("2006-Jan-2_15-04_")

		move_path := path.Dir(online) + "/" + timestamp + path.Base(file.Name())

		fs.Move(online, move_path)
		file.Close()
	}
}

func main() {
	fs.AccessToken = "11.1111111111111111111111111111111111111110.1111111111.1111111111-1111111"
	fs.Client = &http.Client{}

	if len(os.Args) != 3 {
		log.Fatalln("Usage:", os.Args[0], "local", "online")
	}

	RenamePath(os.Args[1], os.Args[2])
}
