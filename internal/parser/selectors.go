package parser

// selectors holds the CSS selectors for a Cambridge dictionary page variant.
type selectors struct {
	headword    string
	phonetic    string
	phoneticUK  string
	phoneticUS  string
	audioUK     string
	audioUS     string
	pos         string
	defBlock    string
	definition  string
	translation string // Vietnamese translation within a def-block or example; empty for English
	example     string
}

var englishSelectors = selectors{
	headword:   ".headword",
	phoneticUK: ".uk .pron",
	phoneticUS: ".us .pron",
	audioUK:    ".uk.dpron-i source",
	audioUS:    ".us.dpron-i source",
	pos:        ".posgram .pos",
	defBlock:   ".def-block",
	definition: ".def",
	example:    ".examp",
}

var englishVietnameseSelectors = selectors{
	headword:    ".english-vietnamese .di-title",
	phonetic:    ".english-vietnamese .pron",
	pos:         ".english-vietnamese .pos",
	defBlock:    ".def-block",
	definition:  ".def",
	translation: ".trans",
	example:     ".examp",
}
