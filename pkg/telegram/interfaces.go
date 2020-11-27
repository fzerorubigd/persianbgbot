package telegram

type (
	button struct {
		keys []string
		text string
	}
	message string
)

func (m message) buttonText() string {
	return string(m)
}

func (b *button) buttonText() string {
	return b.text
}

// Response is a telegram response
type Response interface {
	buttonText() string
}

type Menu interface {
	Reset() Response
	Process(message string) Response
}

// NewButtonResponse create new button list
func NewButtonResponse(text string, items ...string) Response {
	return &button{
		keys: items,
		text: text,
	}
}

// NewTextResponse create a new text message
func NewTextResponse(text string) Response {
	return message(text)
}
