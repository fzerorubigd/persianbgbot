// Package teg contains data for Tiny Epic Galaxies game
package teg

import (
	"embed"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/fzerorubigd/persianbgbot/pkg/menu"
)

//go:embed data
var data embed.FS

// card is a single card in game
type card struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Text        string `yaml:"text"`
	TextPersian string `yaml:"text_persian"`
	Note        string `yaml:"note"`
	NotePersian string `yaml:"note_persian"`
}

func (c *card) Index() string {
	return c.Name
}

func (c *card) Message() string {
	text := fmt.Sprintf(`
<i><b>%s</b></i>
<b>متن اصلی کارت</b>: 
%s
<b>ترجمه</b>: 
%s
`, c.Name, c.Text, c.TextPersian)
	if c.Note != "" {
		text += fmt.Sprintf(`
<b> توضیحات انگلیسی</b>: 
%s
<b>ترجمه</b>: 
%s
`, c.Note, c.NotePersian)
	}
	// Contains all the information
	return text
}

// tinyEpicGalaxies contains cards for bloodrage game
type tinyEpicGalaxies struct {
	types []menu.Item
}

// Caption returns the caption
func (b *tinyEpicGalaxies) Caption() string {
	return "Select one card type:"
}

// Index returns the index
func (b *tinyEpicGalaxies) Index() string {
	return "Tiny Epic Galaxies"
}

// Load return cards inside it
func (b *tinyEpicGalaxies) Load() []menu.Item {
	return b.types
}

func (b *tinyEpicGalaxies) Button() bool {
	return true
}

// loadCards load tiny epic galaxies cards from bin data
func loadCards() (*tinyEpicGalaxies, error) {
	fl, err := data.ReadFile("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result []*card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	var (
		secret []menu.Item
		planet []menu.Item
	)

	for a := range result {
		switch result[a].Type {
		case "Planet":
			planet = append(planet, result[a])
		case "Secret Mission":
			secret = append(secret, result[a])
		default:
			return nil, errors.Errorf("invalid type %q", result[a].Type)
		}
	}

	t := &tinyEpicGalaxies{
		types: []menu.Item{
			menu.NewSimpleNode("Planets", "Choose planet card:", true, planet...),
			menu.NewSimpleNode("Secret Mission", "Choose Secret Mission card:", true, secret...),
			menu.NewSimpleLeaf("About", "Tiny Epic Galaxies Cards\n<b>Translated by</b>: Forud Ghafouri"),
		},
	}

	return t, nil
}

func init() {
	b, err := loadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
