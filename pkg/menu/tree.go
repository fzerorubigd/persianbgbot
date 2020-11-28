package menu

import (
	"errors"
	"sort"
	"strings"

	"github.com/fzerorubigd/persianbgbot/pkg/telegram"
)

// Item is a single item in menu (Node or Leaf)
type Item interface {
	Index() string
}

// Node is an item in menu
type Node interface {
	Item

	Load() []Item
	Caption() string

	Button() bool
}

// Leaf is a leaf in menu
type Leaf interface {
	Item

	Message() string
}

const (
	parentButton  = "⤴️"
	backButton    = "⬅️️"
	forwardButton = "➡️"
	resetButton   = "🤷🏻‍"
)

// memoryMenu is a in memory tree of states, NOT CONCURRENT SAFE
type memoryMenu struct {
	root []Node

	caption    string
	filtered   []Item
	current    []Item
	showButton bool

	start      int
	limit      int
	lastFilter string
}

func (t *memoryMenu) buildFilteredMenu(filter string) bool {
	t.filtered = make([]Item, 0, len(t.current))
	for idx := range t.current {
		if t.current[idx].Index() == filter {
			// exact match
			t.filtered = t.filtered[:0]
			t.filtered = append(t.filtered, t.current[idx])
			return true
		}
		if strings.HasPrefix(t.current[idx].Index(), filter) {
			t.filtered = append(t.filtered, t.current[idx])
		}
	}

	return false
}

func (t *memoryMenu) appendItems(parentItem, backItem, forwardItem bool, items ...string) telegram.Response {
	if !t.showButton {
		return telegram.NewButtonResponse(t.caption, parentButton)
	}
	result := make([]string, 0, len(items)+3)
	if parentItem {
		result = append(result, parentButton)
	}

	result = append(result, items...)

	if backItem {
		result = append(result, backButton)
	}

	if forwardItem {
		result = append(result, forwardButton)
	}

	return telegram.NewButtonResponse(t.caption, result...)
}

func (t *memoryMenu) buildMenu(filter string) telegram.Response {
	if t.start < 0 {
		t.start = 0
	}

	final := t.buildFilteredMenu(filter)
	if len(t.filtered) == 0 {
		return telegram.NewButtonResponse(resetButton, resetButton)
	}

	t.lastFilter = filter

	if len(t.filtered) == 1 && final {
		if node, ok := t.filtered[0].(Node); ok {
			t.root = append(t.root, node)
			t.current = node.Load()
			t.showButton = node.Button()
			t.caption = node.Caption()
			t.start = 0
			return t.buildMenu("")
		}
		txt := t.filtered[0].(Leaf).Message()
		topMenu := t.buildMenu("")
		topMenu.SetText(txt)
		return topMenu
	}

	if len(t.filtered) <= t.limit {
		t.start = 0
		items := make([]string, 0, len(t.filtered))
		for i := range t.filtered {
			items = append(items, t.filtered[i].Index())
		}

		return t.appendItems(len(t.root) > 1, false, false, items...)
	}

	duplicate := make(map[string]bool, len(t.filtered))
	for i := range t.filtered {
		title := t.filtered[i].Index()
		if len(title) > len(filter) {
			title = title[:len(filter)+1]
		}
		duplicate[title] = true
	}

	items := make([]string, 0, len(duplicate))
	for i := range duplicate {
		items = append(items, i)
	}
	sort.Strings(items)
	if len(items) <= t.limit {
		return t.appendItems(len(t.root) > 1, false, false, items...)
	}

	if t.start >= len(items) {
		t.start = 0
	}
	if t.start+t.limit < len(items) {
		return t.appendItems(len(t.root) > 1, t.start > 0, true, items[t.start:t.start+t.limit]...)
	}
	return t.appendItems(len(t.root) > 1, t.start > 0, false, items[t.start:]...)
}

func (t *memoryMenu) Reset() telegram.Response {
	t.root = t.root[:1]
	t.current = t.root[0].Load()
	t.showButton = t.root[0].Button()
	t.caption = t.root[0].Caption()

	return t.buildMenu("")
}

func (t *memoryMenu) Process(message string) telegram.Response {
	if message == resetButton {
		return t.buildMenu("")
	}

	if message == backButton {
		t.start -= t.limit
		return t.buildMenu(t.lastFilter)
	}

	if message == forwardButton {
		t.start += t.limit
		return t.buildMenu(t.lastFilter)
	}

	if message == parentButton {
		if t.lastFilter != "" {
			return t.buildMenu("")
		}
		if len(t.root) > 1 {
			t.root = t.root[:len(t.root)-1]
			t.current = t.root[len(t.root)-1].Load()
			t.showButton = t.root[len(t.root)-1].Button()
			t.caption = t.root[len(t.root)-1].Caption()
			return t.buildMenu("")
		}
	}

	return t.buildMenu(message)
}

type rootMenu struct {
	games []Item
}

func (r *rootMenu) Index() string {
	return ""
}

func (r *rootMenu) Load() []Item {
	return r.games
}

func (r *rootMenu) Caption() string {
	return "Select the game"
}

func (r *rootMenu) Button() bool {
	return true
}

// CreateMemoryMenu creates in memory app
func CreateMemoryMenu(limit int, menu ...Node) (telegram.Menu, error) {
	if len(menu) == 0 {
		return nil, errors.New("at least one root menu is required")
	}

	root := &rootMenu{
		games: make([]Item, 0, len(menu)),
	}
	for n := range menu {
		root.games = append(root.games, menu[n])
	}

	return &memoryMenu{
		root:  []Node{root},
		limit: limit,
	}, nil
}
