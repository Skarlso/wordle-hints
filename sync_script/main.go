package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Hint struct {
	Synonyms    []string `json:"synonyms"`
	Explanation string   `json:"explanation"`
}

func main() {
	// read the local script
	content, err := os.ReadFile(filepath.Join("..", "assets", "js", "words.js"))
	if err != nil {
		log.Fatal(err)
	}

	// trim the javascript declaration
	content = bytes.TrimPrefix(content, []byte("var words = "))
	// parse the slice
	var words []string
	if err := json.Unmarshal(content, &words); err != nil {
		log.Fatal(err)
	}
	words = words[:5]
	// {"word": {"synonyms": [], "explanation": "asdf"}, "anotherword": {}}
	hints := make(map[string]*Hint)

	g := new(errgroup.Group)
	sem := semaphore.NewWeighted(5)

	for _, s := range words {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}

		word := s
		g.Go(func() error {
			defer sem.Release(1)

			resp, err := http.Get(fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word))
			if err != nil {
				if resp.StatusCode == 429 {
					fmt.Println("throttled for word: ", word)
					return nil
				}
				return err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			m := make([]interface{}, 0)
			if err := json.Unmarshal(body, &m); err != nil {
				return err
			}
			if len(m) == 0 {
				return errors.New("explanation is empty")
			}
			hint := &Hint{}
			hint.Explanation = m[0].(map[string]interface{})["meanings"].([]interface{})[0].(map[string]interface{})["definitions"].([]interface{})[0].(map[string]interface{})["definition"].(string)
			hints[word] = hint
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	output, err := json.Marshal(hints)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("hints.json", output, 0766); err != nil {
		log.Fatal(err)
	}
	fmt.Println("all done")
}
