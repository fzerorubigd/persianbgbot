package terraforming

import (
	"embed"
	"fmt"
	"github.com/fzerorubigd/persianbgbot/pkg/menu"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sort"
	"strings"
)

//go:embed data
var data embed.FS

// card is a single card in game
type card struct {
	Name           string `yaml:"Name"`
	Number         string `yaml:"Number"`
	Cost           string `yaml:"Cost"`
	CardType       string `yaml:"CardType"`
	Deck           string `yaml:"Deck"`
	ActionPersian  string `yaml:"Action_Persian"`
	OneTimePersian string `yaml:"OneTime_Persian"`
	ActionEnglish  string `yaml:"Action_English"`
	OneTimeEnglish string `yaml:"OneTime_English"`
}

type byName struct {
	*card
}

type byNumber struct {
	*card
}

func (n *byName) Index() string {
	return n.Name
}

func (n *byNumber) Index() string {
	return n.Number
}

func (c *card) Message() string {
	text := fmt.Sprintf(`
<b>%s: <i>%s</i></b>
<b>هزینه</b>: %s
`, c.CardType, c.Name, c.Cost)
	if c.OneTimeEnglish != "" {
		text += fmt.Sprintf(`<b>متن اصلی اثر یک باره</b>:
%s
<b>ترجمه اثر یکباره</b>
%s
`, c.OneTimeEnglish, c.OneTimePersian)
	}

	if c.ActionEnglish != "" {
		text += fmt.Sprintf(`<b>اکشن انگلیسی</b>:
%s
<b>ترجمه اکشن فارسی</b>
%s
`, c.ActionEnglish, c.ActionPersian)
	}

	text += fmt.Sprintf(`
<u>#%s</u> %s`, c.Number, c.Deck)
	return text
}

type terraformingMars struct {
	items []menu.Item
}

func (t *terraformingMars) Index() string {
	return "Terraforming Mars"
}

func (t *terraformingMars) Load() []menu.Item {
	return t.items
}

func (t *terraformingMars) Caption() string {
	return "Choose the search method: "
}

func (t *terraformingMars) Button() bool {
	return true
}

// loadCards load bloodrage cards from bin data
func loadCards() (*terraformingMars, error) {
	fl, err := data.ReadFile("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result []*card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	var name []menu.Item
	for i := range result {
		name = append(name, &byName{result[i]})
	}
	sort.Slice(name, func(i, j int) bool {
		return strings.Compare(name[i].Index(), name[i].Index()) < 0
	})

	nameSort := menu.NewSimpleNode("By Name", "Search cards by name: ", true, name...)

	var num []menu.Item
	for i := range result {
		if result[i].Number != "" {
			num = append(num, &byNumber{result[i]})
		}
	}
	sort.Slice(num, func(i, j int) bool {
		return strings.Compare(num[i].Index(), num[i].Index()) < 0
	})

	numSort := menu.NewSimpleNode("By Card Number", "Input the card number: ", false, num...)
	tm := &terraformingMars{
		items: []menu.Item{nameSort, numSort},
	}

	tm.items = append(tm.items, menu.NewSimpleLeaf("About", `
<b>Terraforming Mars Cards</b>
<i>Work In Progress</i>
I found the basic card translation in BGG (PDF format) by <b>Ali Nasery</b> (Instagram: @boardgamerboy) 
I couldn't find this ID in instagram, so I have no confirmation from him so far.
Some edit (not much, just few places) and Corporate card translation and converting to YAML by Forud Ghafouri
`))

	return tm, nil
}

func init() {
	b, err := loadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
