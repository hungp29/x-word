package service

import (
	"net/http"
	"sync"

	"github.com/hupham/x-word/internal/model"
)

const workerLimit = 3

type Fetcher interface {
	Fetch(string) (*http.Response, error)
}

type Parser interface {
	Parse(*http.Response) (*model.Word, error)
}

type WordService struct {
	fetcher Fetcher
	parser  Parser
}

func NewWordService(fetcher Fetcher, parser Parser) *WordService {
	return &WordService{fetcher: fetcher, parser: parser}
}

func (s *WordService) GetWord(word string) (*model.Word, error) {
	resp, err := s.fetcher.Fetch(word)
	if err != nil {
		return nil, err
	}
	return s.parser.Parse(resp)
}

func (s *WordService) GetWords(words []string) ([]*model.Word, error) {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make([]*model.Word, 0, len(words))
		firstErr error
	)

	sem := make(chan struct{}, workerLimit)

	for _, word := range words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			resp, err := s.fetcher.Fetch(word)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			parsed, err := s.parser.Parse(resp)
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
