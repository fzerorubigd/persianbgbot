package genericloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/fzerorubigd/persianbgbot/pkg/menu"
	"github.com/pkg/errors"
)

type IndexFileStruct struct {
	Name       string `json:"name"`
	Caption    string `json:"caption"`
	IndexField string `json:"index_field"`
	GroupName  string `json:"group_name"`
	Template   string `json:"template,omitempty"`
}

type GameFileStruct struct {
	Name     string                              `json:"name"`
	Caption  string                              `json:"caption"`
	Indices  []*IndexFileStruct                  `json:"indices"`
	Cards    map[string][]map[string]interface{} `json:"cards"`
	About    string                              `json:"about"`
	Template string                              `json:"template,omitempty"`
}

// Game contains all the cards for a game
type Game struct {
	Name       string
	CationText string
	Items      []menu.Item
	About      string
}

// Caption returns the caption
func (b *Game) Caption() string {
	return b.CationText
}

// Index returns the index
func (b *Game) Index() string {
	return b.Name
}

// Load return cards inside it
func (b *Game) Load() []menu.Item {
	return b.Items
}

func (b *Game) Button() bool {
	// Game is always a button
	return true
}

type Group struct {
	Name        string
	CaptionText string
	Cards       []menu.Item
}

func (a *Group) Caption() string {
	return a.CaptionText
}

func (a *Group) Index() string {
	return a.Name
}

func (a *Group) Load() []menu.Item {
	return a.Cards
}

func (a *Group) Button() bool {
	// Groups are always button
	return true
}

// Card is a single card in game
type Card struct {
	Name string
	Text string
}

func (c *Card) Index() string {
	return c.Name
}

func (c *Card) Message() string {
	return c.Text
}

func buildCard(tpl *template.Template, idx *IndexFileStruct, card map[string]interface{}) (*Card, error) {
	iface, ok := card[idx.IndexField]
	if !ok {
		return nil, errors.Errorf("card has no index field in group: %q", idx.Name)
	}

	index := fmt.Sprint(iface)

	data := struct {
		Card  map[string]interface{}
		Group string
	}{
		Card:  card,
		Group: idx.Name,
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, errors.Wrapf(err, "rendering card %q failed in group %q", index, idx.Name)
	}

	return &Card{
		Name: index,
		Text: buf.String(),
	}, nil
}

func buildGroup(idx *IndexFileStruct, cards []map[string]interface{}) ([]menu.Item, error) {
	result := make([]menu.Item, 0, len(cards))

	tpl, err := template.New(idx.GroupName).Parse(idx.Template)
	if err != nil {
		return nil, errors.Wrapf(err, "template failed for index %q", idx.Name)
	}

	for card := range cards {
		item, err := buildCard(tpl, idx, cards[card])
		if err != nil {
			return nil, errors.Wrap(err, "building card failed")
		}

		result = append(result, item)
	}

	return result, nil
}

func buildIndex(idx *IndexFileStruct, cards map[string][]map[string]interface{}) (*Group, error) {
	result := &Group{
		Name:        idx.Name,
		CaptionText: idx.Caption,
	}

	for group := range cards {
		if idx.GroupName != "" && idx.GroupName != group {
			continue
		}

		cards, err := buildGroup(idx, cards[group])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build group %q", group)
		}
		result.Cards = append(result.Cards, cards...)
	}

	return result, nil
}

// loadCards load card from an io.Reader
func loadCards(r io.Reader) (*Game, error) {
	var game GameFileStruct
	if err := json.NewDecoder(r).Decode(&game); err != nil {
		return nil, errors.Wrap(err, "JSON decode failed")
	}

	result := &Game{
		Name:       game.Name,
		CationText: game.Caption,
		About:      game.About,
	}

	for idx := range game.Indices {
		if game.Indices[idx].Template == "" {
			// Fallback to the game template
			game.Indices[idx].Template = game.Template
		}
		index, err := buildIndex(game.Indices[idx], game.Cards)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build index %q", game.Indices[idx].Name)
		}
		result.Items = append(result.Items, index)
	}

	result.Items = append(result.Items, menu.NewSimpleLeaf("About", game.About))

	return result, nil
}

// RegisterCard read the cards from an ioReader and register them as a menu
func RegisterCard(r io.Reader) error {
	game, err := loadCards(r)
	if err != nil {
		return errors.Wrap(err, "loading game cards failed")
	}

	menu.RegisterGame(game)
	return nil
}
