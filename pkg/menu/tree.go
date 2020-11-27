package menu

import (
	"errors"
	"sort"
	"strings"

	"github.com/fzerorubigd/persianbgbot/pkg/telegram"
)

type Item interface {
	Index() string
}

type Node interface {
	Item

	Load() []Item
	Caption() string
}

type Leaf interface {
	Item

	Message() string
}

const (
	parent  = "⤴️"
	back    = "⬅️️"
	forward = "➡️"
	reset   = "🤷🏻‍"
)

// MemTree is a in memory tree of states, NOT CONCURRENT SAFE
type MemTree struct {
	root []Node

	caption  string
	filtered []Item
	current  []Item

	start      int
	limit      int
	lastFilter string
}

func (t *MemTree) buildFilteredMenu(filter string) bool {
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

func appendItems(parentItem, backItem, forwardItem bool, text string, items ...string) telegram.Response {
	result := make([]string, 0, len(items)+3)
	if parentItem {
		result = append(result, parent)
	}

	result = append(result, items...)

	if backItem {
		result = append(result, back)
	}

	if forwardItem {
		result = append(result, forward)
	}

	return telegram.NewButtonResponse(text, result...)
}

func (t *MemTree) buildMenu(filter string) telegram.Response {
	if t.start < 0 {
		t.start = 0
	}

	final := t.buildFilteredMenu(filter)
	if len(t.filtered) == 0 {
		return telegram.NewButtonResponse(reset, reset)
	}

	t.lastFilter = filter

	if len(t.filtered) == 1 && final {
		if node, ok := t.filtered[0].(Node); ok {
			t.root = append(t.root, node)
			t.current = node.Load()
			t.caption = node.Caption()
			t.start = 0
			return t.buildMenu("")
		}

		return telegram.NewTextResponse(t.filtered[0].(Leaf).Message())
	}

	if len(t.filtered) <= t.limit {
		t.start = 0
		items := make([]string, 0, len(t.filtered))
		for i := range t.filtered {
			items = append(items, t.filtered[i].Index())
		}

		return appendItems(len(t.root) > 1, false, false, t.caption, items...)
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
		return appendItems(len(t.root) > 1, false, false, t.caption, items...)
	}

	if t.start >= len(items) {
		t.start = 0
	}
	if t.start+t.limit < len(items) {
		return appendItems(len(t.root) > 1, t.start > 0, true, t.caption, items[t.start:t.start+t.limit]...)
	}
	return appendItems(len(t.root) > 1, t.start > 0, false, t.caption, items[t.start:]...)
}

func (t *MemTree) Reset() telegram.Response {
	t.root = t.root[:1]
	t.current = t.root[0].Load()
	t.caption = t.root[0].Caption()

	return t.buildMenu("")
}

func (t *MemTree) Process(message string) telegram.Response {
	if message == reset {
		return t.buildMenu("")
	}

	if message == back {
		t.start -= t.limit
		return t.buildMenu(t.lastFilter)
	}

	if message == forward {
		t.start += t.limit
		return t.buildMenu(t.lastFilter)
	}

	if message == parent {
		if t.lastFilter != "" {
			return t.buildMenu("")
		}
		if len(t.root) > 1 {
			t.root = t.root[:len(t.root)-1]
			t.current = t.root[len(t.root)-1].Load()
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

func CreateMenu(limit int, menu ...Node) (telegram.Menu, error) {
	if len(menu) == 0 {
		return nil, errors.New("at least one root menu is required")
	}

	root := &rootMenu{
		games: make([]Item, 0, len(menu)),
	}
	for n := range menu {
		root.games = append(root.games, menu[n])
	}

	return &MemTree{
		root:  []Node{root},
		limit: limit,
	}, nil
}
