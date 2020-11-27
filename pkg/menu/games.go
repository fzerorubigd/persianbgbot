package menu

import (
	"sort"
	"strings"
	"sync"
)

var (
	games = sync.Map{}
)

// RegisterGame register a game into the list
func RegisterGame(game Node) {
	games.Store(game.Index(), game)
}

// AllGames return all registered games
func AllGames() []Node {
	var result []Node
	games.Range(func(key, value interface{}) bool {
		result = append(result, value.(Node))

		return true
	})

	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Index(), result[i].Index()) < 0
	})

	return result
}
