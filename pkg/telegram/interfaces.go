package telegram

type (
	button struct {
		keys []string
		text string
	}
	message struct {
		text string
	}
)

func (m *message) buttonText() string {
	return m.text
}

func (m *message) SetText(s string) {
	m.text = s
}

func (b *button) buttonText() string {
	return b.text
}

func (b *button) SetText(s string) {
	b.text = s
}

// Response is a telegram response
type Response interface {
	buttonText() string

	SetText(string)
}

// Menu is the menu to handle the all menus in bot
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
	return &message{
		text: text,
	}
}
