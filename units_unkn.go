package gomorphy

type unknAnalyzer struct {
	name     string
	selfTag  string
	terminal bool
}

func newUnknAnalyzer() *unknAnalyzer {
	return &unknAnalyzer{
		name:     "UnknAnalyzer",
		selfTag:  "UNKN",
		terminal: true,
	}
}

func (self *unknAnalyzer) Name() string {
	return self.name
}

func (self *unknAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *unknAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	if len(seenParses) != 0 {
		return []*Parse{}
	}
	parse := newParse(wordLower, self.selfTag, 1.0)
	parse.NormalForm = wordLower
	parse.MethodsStack = []Method{{Analyzer: self, WordOrStack: word}}
	return []*Parse{parse}
}

func (self *unknAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if len(seenTags) != 0 {
		return nil
	}

	return []*opencorporaTag{newOpencorporaTag(self.selfTag)}
}

func (self *unknAnalyzer) getLexeme(form *Parse) []*Parse {
	return []*Parse{form}
}

func (self *unknAnalyzer) normalized(form *Parse) *Parse {
	return form
}
