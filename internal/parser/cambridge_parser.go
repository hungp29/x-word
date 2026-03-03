package parser

import (
	"net/http"
	"strings"

	"github.com/hupham/x-word/internal/model"

	"github.com/PuerkitoBio/goquery"
)

type CambridgeParser struct{}

func NewCambridgeParser() *CambridgeParser {
	return &CambridgeParser{}
}

func (p *CambridgeParser) Parse(resp *http.Response) (*model.Word, error) {

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	word := &model.Word{}

	word.Text = doc.Find(".headword").First().Text()
	word.PhoneticUK = doc.Find(".uk .pron").First().Text()
	word.PhoneticUS = doc.Find(".us .pron").First().Text()

	// Get audio URLs
	ukSrc, exists := doc.Find(".uk.dpron-i source").First().Attr("src")
	if exists {
		word.AudioUK = "https://dictionary.cambridge.org" + strings.TrimSpace(ukSrc)
	}
	usSrc, exists := doc.Find(".us.dpron-i source").First().Attr("src")
	if exists {
		word.AudioUS = "https://dictionary.cambridge.org" + strings.TrimSpace(usSrc)
	}

	// Get all unique parts of speech
	posMap := make(map[string]struct{})
	doc.Find(".posgram .pos").Each(func(i int, s *goquery.Selection) {
		pos := s.Text()
		if pos != "" {
			posMap[pos] = struct{}{}
		}
	})

	for p := range posMap {
		word.PartOfSpeech = append(word.PartOfSpeech, p)
	}

	// Get all meanings
	doc.Find(".def-block").Each(func(i int, s *goquery.Selection) {
		def := s.Find(".def").Text()

		var examples []string
		s.Find(".examp").Each(func(j int, e *goquery.Selection) {
			examples = append(examples, strings.TrimSpace(e.Text()))
		})

		word.Meanings = append(word.Meanings, model.Meaning{
			Definition: def,
			Examples:   examples,
		})
	})

	return word, nil
}
