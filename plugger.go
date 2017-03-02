package main

import (
	"log"
	"os"
	"path/filepath"
	"plugin"
)

type pastebin interface {
	Get() (string, error)
	Put(string) error
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	paths, err := filepath.Glob(filepath.Join(dir, "*.so"))
	if err != nil {
		return err
	}

	for _, path := range paths {
		if p, err := plugin.Open(path); err != nil {
			return err
		}
	}

}
