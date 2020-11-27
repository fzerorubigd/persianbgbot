package bloodrage

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/fzerorubigd/persianbgbot/pkg/menu"
)

// Card is a single card in game
type Card struct {
	Name        string `yaml:"Name"`
	STR         string `yaml:"STR"`
	PlayerCount string `yaml:"PlayerCount"`
	English     string `yaml:"English"`
	Persian     string `yaml:"Persian"`
}

func (c *Card) Index() string {
	return c.Name
}

func (c *Card) Message() string {
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

type Age struct {
	Name  string
	Cards []menu.Item
}

func (a *Age) Caption() string {
	return "Select the card to show: "
}

func (a *Age) Index() string {
	return a.Name
}

func (a *Age) Load() []menu.Item {
	return a.Cards
}

type BloodRage struct {
	Ages []menu.Item
}

func (b *BloodRage) Caption() string {
	return "Select one bloodrage Age:"
}

func (b *BloodRage) Index() string {
	return "Bloodrage"
}

func (b *BloodRage) Load() []menu.Item {
	return b.Ages
}

// LoadCards load bloodrage cards from bin data
func LoadCards() (*BloodRage, error) {
	fl, err := Asset("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result map[string][]*Card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	b := &BloodRage{}
	for a := range result {
		age := &Age{
			Name:  a,
			Cards: nil,
		}
		for _, card := range result[a] {
			age.Cards = append(age.Cards, card)
		}
		b.Ages = append(b.Ages, age)
	}

	b.Ages = append(b.Ages, menu.NewSimpleLeaf("About","Bloodrage Cards\n<b>Translated by</b>: Forud Ghafouri"))

	return b, nil
}

func init() {
	b, err := LoadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
