package parser

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hupham/x-word/internal/model"
)

// CambridgeParser parses the English-only Cambridge dictionary page.
type CambridgeParser struct {
	sel selectors
}

func NewCambridgeParser() *CambridgeParser {
	return &CambridgeParser{sel: englishSelectors}
}

func (p *CambridgeParser) Parse(resp *http.Response) (*model.Word, error) {
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	word := &model.Word{}
	word.Text = doc.Find(p.sel.headword).First().Text()
	word.PhoneticUK = doc.Find(p.sel.phoneticUK).First().Text()
	word.PhoneticUS = doc.Find(p.sel.phoneticUS).First().Text()

	if src, exists := doc.Find(p.sel.audioUK).First().Attr("src"); exists {
		word.AudioUK = "https://dictionary.cambridge.org" + strings.TrimSpace(src)
	}
	if src, exists := doc.Find(p.sel.audioUS).First().Attr("src"); exists {
		word.AudioUS = "https://dictionary.cambridge.org" + strings.TrimSpace(src)
	}

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
		var examples []string
		s.Find(p.sel.example).Each(func(_ int, e *goquery.Selection) {
			examples = append(examples, strings.TrimSpace(e.Text()))
		})
		word.Meanings = append(word.Meanings, model.Meaning{
			Definition: strings.TrimSpace(s.Find(p.sel.definition).Text()),
			Examples:   examples,
		})
	})

	return word, nil
}
