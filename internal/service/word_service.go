package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/hupham/x-word/internal/model"
)

const workerLimit = 3

// Dictionary identifies a Cambridge dictionary variant and owns its URL template.
type Dictionary string

const (
	DictionaryEnglish           Dictionary = "english"
	DictionaryEnglishVietnamese Dictionary = "english-vietnamese"
)

// DictionaryURLTemplates maps each supported dictionary to its URL template.
// Exported so callers (e.g. handlers) can validate user-supplied dictionary names.
var DictionaryURLTemplates = map[Dictionary]string{
	DictionaryEnglish:           "https://dictionary.cambridge.org/dictionary/english/%s",
	DictionaryEnglishVietnamese: "https://dictionary.cambridge.org/vi/dictionary/english-vietnamese/%s",
}

// wordURL returns the lookup URL for the given word, or an error if the dictionary is unknown.
func (d Dictionary) wordURL(word string) (string, error) {
	tmpl, ok := DictionaryURLTemplates[d]
	if !ok {
		return "", fmt.Errorf("unknown dictionary %q", d)
	}
	return fmt.Sprintf(tmpl, word), nil
}

type Fetcher interface {
	Fetch(url string) (*http.Response, error)
}

type Parser interface {
	Parse(*http.Response) (*model.Word, error)
}

type WordService struct {
	fetcher Fetcher
	parsers map[Dictionary]Parser
}

func NewWordService(fetcher Fetcher, parsers map[Dictionary]Parser) *WordService {
	return &WordService{fetcher: fetcher, parsers: parsers}
}

func (s *WordService) parserFor(dict Dictionary) (Parser, error) {
	p, ok := s.parsers[dict]
	if !ok {
		return nil, fmt.Errorf("no parser registered for dictionary %q", dict)
	}
	return p, nil
}

func (s *WordService) GetWord(word string, dict Dictionary) (*model.Word, error) {
	p, err := s.parserFor(dict)
	if err != nil {
		return nil, err
	}
	url, err := dict.wordURL(word)
	if err != nil {
		return nil, err
	}
	resp, err := s.fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	return p.Parse(resp)
}

func (s *WordService) GetWords(words []string, dict Dictionary) ([]*model.Word, error) {
	p, err := s.parserFor(dict)
	if err != nil {
		return nil, err
	}

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		results  = make([]*model.Word, 0, len(words))
		firstErr error
	)

	sem := make(chan struct{}, workerLimit)

	for _, word := range words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			url, err := dict.wordURL(word)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			resp, err := s.fetcher.Fetch(url)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			parsed, err := p.Parse(resp)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			results = append(results, parsed)
			mu.Unlock()
		}(word)
	}

	wg.Wait()
	return results, firstErr
}
