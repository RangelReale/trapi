package trapi

type SourceStackItemType int

const (
	SITEM_NONE SourceStackItemType = iota
	SITEM_DATATYPE
	SITEM_API
	SITEM_TEXT
)

func (s SourceStackItemType) String() string {
	switch s {
	case SITEM_NONE:
		return "SITEM_NONE"
	case SITEM_DATATYPE:
		return "SITEM_DATATYPE"
	case SITEM_API:
		return "SITEM_API"
	case SITEM_TEXT:
		return "SITEM_TEXT"
	}
	return "SITEM_UNKNOWN"
}

type SourceStackData struct {
	ItemType      SourceParseItemType
	StackItemType SourceStackItemType
	Item          interface{}
	StackItem     interface{}
}

type SourceParseStack struct {
	s []*SourceStackData
}

func NewSourceParseStack() *SourceParseStack {
	return &SourceParseStack{make([]*SourceStackData, 0)}
}

func (s *SourceParseStack) Push(v *SourceStackData) {
	s.s = append(s.s, v)
}

func (s *SourceParseStack) Len() int {
	return len(s.s)
}

func (s *SourceParseStack) Top() *SourceStackData {
	l := len(s.s)
	if l == 0 {
		return nil
	}

	return s.s[l-1]
}

func (s *SourceParseStack) Pop() *SourceStackData {
	l := len(s.s)
	if l == 0 {
		return nil
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res
}
