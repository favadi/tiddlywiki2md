package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	f, err := os.Create(filepath.Join(out, n.fileName()))
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write([]byte(n.Title + "\n")); err != nil {
		return err
	}
	if _, err = f.Write([]byte(strings.Repeat("#", len(n.Title)) + "\n")); err != nil {
		return err
	}
	if _, err = f.Write([]byte("\n")); err != nil {
		return err
	}

	if _, err = f.Write([]byte(n.Text)); err != nil {
		return err
	}

	return f.Close()
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
