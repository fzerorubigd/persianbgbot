package menu

type simpleMenuItem struct {
	text  string
	index string
}

type simpleNode struct {
	index   string
	caption string
	leaf    []Item
}

func (s *simpleNode) Index() string {
	return s.index
}

func (s *simpleNode) Load() []Item {
	return s.leaf
}

func (s *simpleNode) Caption() string {
	return s.caption
}

func (s *simpleMenuItem) Index() string {
	return s.index
}

func (s *simpleMenuItem) Message() string {
	return s.text
}

// NewSimpleLeaf is the simple menu item
func NewSimpleLeaf(index, text string) Leaf {
	return &simpleMenuItem{
		text:  text,
		index: index,
	}
}

// NewSimpleNode is for creating simple node
func NewSimpleNode(index, caption string, items ...Item) Node {
	return &simpleNode{
		index:   index,
		caption: caption,
		leaf:    items,
	}
}
