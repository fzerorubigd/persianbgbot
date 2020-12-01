package tinytown

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
	Name    string `yaml:"Name"`
	English string `yaml:"English"`
	Persian string `yaml:"Persian"`
}

func (c *card) Index() string {
	return c.Name
}

func (c *card) Message() string {
	text := fmt.Sprintf(`
<i><b>%s</b></i>
<b>متن اصلی</b>: 
%s
<b>ترجمه</b>: 
%s
`, c.Name, c.English, c.Persian)
	// Contains all the information
	return text
}

type cardType struct {
	Name  string
	Cards []menu.Item
}

func (a *cardType) Caption() string {
	return "Select the card to show: "
}

func (a *cardType) Index() string {
	return a.Name
}

func (a *cardType) Load() []menu.Item {
	return a.Cards
}

func (a *cardType) Button() bool {
	return true
}

// tinyTown contains cards for Tiny Town game
type tinyTown struct {
	types []menu.Item
}

// Caption returns the caption
func (b *tinyTown) Caption() string {
	return "Select one Tiny Town cardType:"
}

// Index returns the index
func (b *tinyTown) Index() string {
	return "Tiny Town"
}

// Load return cards inside it
func (b *tinyTown) Load() []menu.Item {
	return b.types
}

func (b *tinyTown) Button() bool {
	return true
}

// loadCards load bloodrage cards from bin data
func loadCards() (*tinyTown, error) {
	fl, err := data.ReadFile("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result map[string][]*card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	b := &tinyTown{}
	for a := range result {
		typ := &cardType{
			Name:  a,
			Cards: nil,
		}
		for _, card := range result[a] {
			typ.Cards = append(typ.Cards, card)
		}
		b.types = append(b.types, typ)
	}

	b.types = append(b.types, menu.NewSimpleLeaf("About", "Tiny Town Cards\nOnly Monument cards\n<b>Translated by</b>: Forud Ghafouri"))

	return b, nil
}

func init() {
	b, err := loadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
