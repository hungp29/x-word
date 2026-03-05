package grpcserver

import (
	wordv1 "github.com/hungp29/x-proto/gen/go/word/v1"
	"github.com/hungp29/x-word/internal/model"
	"github.com/hungp29/x-word/internal/service"
)

// serviceDictFromProto maps proto Dictionary to service.Dictionary.
func serviceDictFromProto(d wordv1.Dictionary) service.Dictionary {
	switch d {
	case wordv1.Dictionary_DICTIONARY_ENGLISH_VIETNAMESE:
		return service.DictionaryEnglishVietnamese
	case wordv1.Dictionary_DICTIONARY_ENGLISH, wordv1.Dictionary_DICTIONARY_UNSPECIFIED:
	default:
	}
	return service.DictionaryEnglish
}

// wordToProto converts model.Word to wordv1.Word.
func wordToProto(w *model.Word) *wordv1.Word {
	if w == nil {
		return nil
	}
	meanings := make([]*wordv1.Meaning, len(w.Meanings))
	for i, m := range w.Meanings {
		meanings[i] = &wordv1.Meaning{
			Definition: m.Definition,
			Examples:   m.Examples,
		}
	}
	return &wordv1.Word{
		Text:         w.Text,
		Phonetic:     w.Phonetic,
		PhoneticUk:   w.PhoneticUK,
		PhoneticUs:   w.PhoneticUS,
		AudioUk:      w.AudioUK,
		AudioUs:      w.AudioUS,
		PartOfSpeech: w.PartOfSpeech,
		Meanings:     meanings,
	}
}
