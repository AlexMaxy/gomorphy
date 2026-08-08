package gomorphy

type analyzer interface {
	Name() string
	isTerminal() bool
	parse(string, string, map[seenParse]bool) []*Parse
	tag(string, string, map[string]bool) []*opencorporaTag
	getLexeme(*Parse) []*Parse
	normalized(*Parse) *Parse
}

func (m *MorphAnalyzer) addAnalyzer(a analyzer) {
	m.analyzers = append(m.analyzers, a)
}
