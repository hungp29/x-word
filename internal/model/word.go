package model

type Word struct {
	Text         string    `json:"text"`
	Phonetic     string    `json:"phonetic"`
	PhoneticUK   string    `json:"phonetic_uk"`
	PhoneticUS   string    `json:"phonetic_us"`
	AudioUK      string    `json:"audio_uk"`
	AudioUS      string    `json:"audio_us"`
	PartOfSpeech []string  `json:"part_of_speech"`
	Meanings     []Meaning `json:"meanings"`
}

type Meaning struct {
	Definition string   `json:"definition"`
	Examples   []string `json:"examples"`
}
