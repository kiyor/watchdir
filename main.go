/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 07-25-2016

* Last Modified : Wed 14 Sep 2016 04:51:09 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

var (
	dir   = flag.String("d", "./", "monitor dir")
	chAdd = make(chan string)
	chDel = make(chan string)

	version = flag.Bool("version", false, "output version and exit")

	buildtime string
	VER       = "1.0"
)

func init() {
	flag.Parse()
	if *version {
		fmt.Printf("%s.%s\n", VER, buildtime)
		os.Exit(0)
	}
}

func main() {
	done := make(chan bool)
	go watch(chAdd, chDel)
	err := filepath.Walk(*dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err.Error())
			// 			return err
		} else if f.IsDir() {
			chAdd <- path
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	<-done
}

func watch(chAdd, chDel chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op == fsnotify.Create {
					f, err := os.Stat(event.Name)
					if err != nil {
						log.Println("error", err.Error())
						continue
					}
					if f.IsDir() {
						chAdd <- event.Name
					}
				}
				if event.Op == fsnotify.Remove {
					chDel <- event.Name
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		for d := range chAdd {
			err = watcher.Add(d)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("monitor", d)
		}
	}()
	for d := range chDel {
		watcher.Remove(d)
	}
}
