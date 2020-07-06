package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Note struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (n *Note) fileName() string {
	var fileName = strings.ReplaceAll(n.Title, "/", "")
	return fmt.Sprintf("%s.md", fileName)
}

func (n *Note) isSkipped() bool {
	return false
}

func (n *Note) Write(out string) error {
	if n.isSkipped() {
		return nil
	}

	return ioutil.WriteFile(filepath.Join(out, n.fileName()), []byte(n.Text), 0644)
}

func main() {
	var (
		input  = flag.String("input", "tiddlers.json", "TiddlyWiki JSON export file")
		output = flag.String("output", "output", "output directory")

		err error
	)

	flag.Parse()

	if err = os.MkdirAll(*output, 0755); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	var notes []Note
	if err = json.NewDecoder(f).Decode(&notes); err != nil {
		log.Fatal(err)
	}

	for _, note := range notes {
		if err = note.Write(*output); err != nil {
			log.Fatal(err)
		}
	}
}
