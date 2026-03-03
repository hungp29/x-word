package service

import (
	"net/http"
	"sync"

	"github.com/hupham/x-word/internal/model"
)

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
	var wg sync.WaitGroup
	resultChan := make(chan *model.Word, len(words))
	errorChan := make(chan error, len(words))

	workerLimit := 3
	sem := make(chan struct{}, workerLimit)

	for _, word := range words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			resp, err := s.fetcher.Fetch(word)
			if err != nil {
				errorChan <- err
				return
			}

			parsedWord, err := s.parser.Parse(resp)
			if err != nil {
				errorChan <- err
				return
			}

			resultChan <- parsedWord
		}(word)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	results := make([]*model.Word, 0, len(words))
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}
