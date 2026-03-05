package parser

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hungp29/x-word/internal/model"
)

// EnglishVietnameseParser parses the Cambridge English-Vietnamese dictionary page.
type EnglishVietnameseParser struct {
	sel selectors
}

func NewEnglishVietnameseParser() *EnglishVietnameseParser {
	return &EnglishVietnameseParser{sel: englishVietnameseSelectors}
}

func (p *EnglishVietnameseParser) Parse(resp *http.Response) (*model.Word, error) {
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	word := &model.Word{}
	word.Text = doc.Find(p.sel.headword).First().Text()
	word.Phonetic = doc.Find(p.sel.phonetic).First().Text()

	posMap := make(map[string]struct{})
	doc.Find(p.sel.pos).Each(func(_ int, s *goquery.Selection) {
		if pos := strings.TrimSpace(s.Text()); pos != "" {
			posMap[pos] = struct{}{}
		}
	})
	for pos := range posMap {
		word.PartOfSpeech = append(word.PartOfSpeech, pos)
	}

	doc.Find(p.sel.defBlock).Each(func(_ int, s *goquery.Selection) {
		// Prefer the Vietnamese translation; fall back to the English gloss.
		definition := strings.TrimSpace(s.Find(p.sel.translation).First().Text())
		if definition == "" {
			definition = strings.TrimSpace(s.Find(p.sel.definition).Text())
		}

		var examples []string
		s.Find(p.sel.example).Each(func(_ int, e *goquery.Selection) {
			examples = append(examples, strings.TrimSpace(e.Text()))
		})

		word.Meanings = append(word.Meanings, model.Meaning{
			Definition: definition,
			Examples:   examples,
		})
	})

	return word, nil
}
