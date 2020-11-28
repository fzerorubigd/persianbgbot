package terraforming

import (
	"fmt"
	"github.com/fzerorubigd/persianbgbot/pkg/menu"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sort"
	"strings"
)

// card is a single card in game
type card struct {
	Name           string `yaml:"Name"`
	Number         int    `yaml:"Number"`
	Cost           int    `yaml:"Cost"`
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
	return fmt.Sprint(n.Number)
}

func (c *card) Message() string {
	text := fmt.Sprintf(`
<b>%s: <i>%s</i></b>
<b>هزینه</b>: %d
`, c.CardType, c.Name, c.Cost)
	if c.OneTimeEnglish != "" {
		text += fmt.Sprintf(`<b>متن اصلی اثر یک باره</b>:
%s
<b>ترجمه اثر یکباره</b>
%s
`, c.OneTimeEnglish, c.OneTimeEnglish)
	}

	if c.ActionEnglish != "" {
		text += fmt.Sprintf(`<b>اکشن انگلیسی</b>:
%s
<b>ترجمه اکشن فارسی</b>
%s
`, c.ActionEnglish, c.ActionPersian)
	}

	text += fmt.Sprintf(`
<u>#%d</u> %s`, c.Number, c.Deck)
	return text
}

type customIndexer struct {
	index   string
	caption string
	items   []menu.Item
}

func (c *customIndexer) Index() string {
	return c.index
}

func (c *customIndexer) Load() []menu.Item {
	return c.items
}

func (c *customIndexer) Caption() string {
	return c.caption
}

type terraformingMars struct {
	items []menu.Item
}

func (t terraformingMars) Index() string {
	return "Terraforming Mars"
}

func (t terraformingMars) Load() []menu.Item {
	return t.items
}

func (t terraformingMars) Caption() string {
	return "Choose the search method: "
}

// loadCards load bloodrage cards from bin data
func loadCards() (*terraformingMars, error) {
	fl, err := Asset("data/cards.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "can not load the files from bin storage")
	}
	var result []*card
	if err := yaml.Unmarshal(fl, &result); err != nil {
		return nil, errors.Wrap(err, "load yaml data failed")
	}

	nameSort := customIndexer{
		index:   "By Name",
		caption: "Search cards by name: ",
		items:   nil,
	}
	for i := range result {
		nameSort.items = append(nameSort.items, &byName{result[i]})
	}
	sort.Slice(nameSort.items, func(i, j int) bool {
		return strings.Compare(nameSort.items[i].Index(), nameSort.items[i].Index()) < 0
	})

	numSort := customIndexer{
		index:   "By Card Number",
		caption: "Search cards by number (You can type number): ",
		items:   nil,
	}
	for i := range result {
		numSort.items = append(numSort.items, &byNumber{result[i]})
	}
	sort.Slice(numSort.items, func(i, j int) bool {
		return strings.Compare(numSort.items[i].Index(), numSort.items[i].Index()) < 0
	})

	return &terraformingMars{
		items: []menu.Item{&nameSort, &numSort},
	}, nil
}

func init() {
	b, err := loadCards()
	if err != nil {
		panic("invalid data")
	}
	menu.RegisterGame(b)
}
