package parser

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hupham/x-word/internal/model"
)

// EnglishVietnameseParser parses the Cambridge English-Vietnamese dictionary page.
// Definitions are Vietnamese translations extracted from .trans elements;
// examples are taken from .examp (English sentence) paired with .trans inside it.
type EnglishVietnameseParser struct{}

func NewEnglishVietnameseParser() *EnglishVietnameseParser {
	return &EnglishVietnameseParser{}
}

func (p *EnglishVietnameseParser) Parse(resp *http.Response) (*model.Word, error) {
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	word := &model.Word{}
	word.Text = doc.Find(".english-vietnamese .di-title").First().Text()
	word.Phonetic = doc.Find(".english-vietnamese .pron").First().Text()

	posMap := make(map[string]struct{})
	doc.Find(".english-vietnamese .pos").Each(func(_ int, s *goquery.Selection) {
		if pos := strings.TrimSpace(s.Text()); pos != "" {
			posMap[pos] = struct{}{}
		}
	})
	for pos := range posMap {
		word.PartOfSpeech = append(word.PartOfSpeech, pos)
	}

	doc.Find(".def-block").Each(func(_ int, s *goquery.Selection) {
		// Prefer the Vietnamese translation; fall back to the English gloss.
		translation := strings.TrimSpace(s.Find(".trans").First().Text())
		if translation == "" {
			translation = strings.TrimSpace(s.Find(".def").Text())
		}

		var examples []string
		s.Find(".examp").Each(func(_ int, e *goquery.Selection) {
			examples = append(examples, strings.TrimSpace(e.Text()))
		})

		word.Meanings = append(word.Meanings, model.Meaning{
			Definition: translation,
			Examples:   examples,
		})
	})

	return word, nil
}
