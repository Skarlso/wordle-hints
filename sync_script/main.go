package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"
)

type Hint struct {
	Synonym     string `json:"synonym"`
	Explanation string `json:"explanation"`
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
	// words = words[:500]
	// {"word": {"synonyms": [], "explanation": "asdf"}, "anotherword": {}}
	hints := make(map[string]*Hint)

	var word string
	for len(words) > 0 {
		word, words = words[0], words[1:]
		fmt.Println("at word: ", word)
		hint := &Hint{}
		entryURL := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
		synURL := fmt.Sprintf("https://api.datamuse.com/words?ml=%s", word)
		g := new(errgroup.Group)
		g.Go(func() error {
			resp, err := http.Get(entryURL)
			if err != nil {
				fmt.Println("failed entry for word: ", word)
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode == 429 {
				return errors.New("retry")
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			m := make([]interface{}, 0)
			if err := json.Unmarshal(body, &m); err != nil {
				fmt.Println("entry failed to unmarshal")
				fmt.Println("body: ", string(body))
				fmt.Println("word: ", word)
				hint.Explanation = "no explanation, sorry"
				return nil
			}
			if len(m) == 0 {
				return err
			}
			hint.Explanation = m[0].(map[string]interface{})["meanings"].([]interface{})[0].(map[string]interface{})["definitions"].([]interface{})[0].(map[string]interface{})["definition"].(string)
			return nil
		})
		g.Go(func() error {
			resp, err := http.Get(synURL)
			if err != nil {
				fmt.Println("failed syn for word: ", word)
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode == 429 {
				return errors.New("retry")
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			m := make([]interface{}, 0)
			if err := json.Unmarshal(body, &m); err != nil {
				fmt.Println("syn failed to unmarshal")
				fmt.Println("body: ", string(body))
				fmt.Println("word: ", word)
				hint.Synonym = "no synonym, sorry"
				return nil
			}
			if len(m) == 0 {
				return err
			}
			hint.Synonym = m[0].(map[string]interface{})["word"].(string)
			return nil
		})

		if err := g.Wait(); err != nil {
			if err.Error() == "retry" {
				time.Sleep(5 * time.Minute)
				// re-add the word at the end to retry
				words = append(words, word)
				continue
			}
			log.Fatal(err)
		}

		hints[word] = hint
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
