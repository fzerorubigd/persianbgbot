// Package bloodrage contains data for blodrage game
package bloodrage

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/fzerorubigd/persianbgbot/pkg/menu"
)

// card is a single card in game
type card struct {
	Name        string `yaml:"Name"`
	STR         string `yaml:"STR"`
	PlayerCount string `yaml:"PlayerCount"`
	English     string `yaml:"English"`
	Persian     string `yaml:"Persian"`
}

func (c *card) Index() string {
	return c.Name
}

func (c *card) Message() string {
	text := fmt.Sprintf(`
<i><b>%s</b></i>
<b>قدرت کارت (STR)</b>: 
%s
<b>متن اصلی</b>: 
%s
<b>ترجمه</b>: 
%s
`, c.Name, c.STR, c.English, c.Persian)
	// Contains all the information
	return text
}

type age struct {
	Name  string
	Cards []menu.Item
}

func (a *age) Caption() string {
	return "Select the card to show: "
}

func (a *age) Index() string {
	return a.Name
}

func (a *age) Load() []menu.Item {
	return a.Cards
}

func (a *age) Button() bool {
	return true
}

// bloodRage contains cards for bloodrage game
type bloodRage struct {
	Ages []menu.Item
}

// Caption returns the caption
func (b *bloodRage) Caption() string {
	return "Select one bloodrage age:"
}

// Index returns the index
func (b *bloodRage) Index() string {
	return "Bloodrage"
}

// Load return cards inside it
func (b *bloodRage) Load() []menu.Item {
	return b.Ages
}

func (b *bloodRage) Button() bool {
	return true
}

// loadCards load bloodrage cards from bin data
func loadCards() (*bloodRage, error) {
	fl, err := Asset("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result map[string][]*card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	b := &bloodRage{}
	for a := range result {
		age := &age{
			Name:  a,
			Cards: nil,
		}
		for _, card := range result[a] {
			age.Cards = append(age.Cards, card)
		}
		b.Ages = append(b.Ages, age)
	}

	b.Ages = append(b.Ages, menu.NewSimpleLeaf("About", "Bloodrage Cards\n<b>Translated by</b>: Forud Ghafouri"))

	return b, nil
}

func init() {
	b, err := loadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
