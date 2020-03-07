package domain

type SpeechRequest struct {
	Text         string
	LanguageCode string
	SsmlGender   string
	VoiceName    string
}

type SpeechExampleError struct {
	Message string
}
