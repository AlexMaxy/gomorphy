package gomorphy

type dictionaryAnalyzer struct {
	name     string
	terminal bool
	*wordsDawg
}

func newDictionaryAnalyzer(dawg *wordsDawg) *dictionaryAnalyzer {
	return &dictionaryAnalyzer{
		name:      "DictionaryAnalyzer",
		terminal:  false,
		wordsDawg: dawg,
	}
}

func (self *dictionaryAnalyzer) Name() string {
	return self.name
}

func (self *dictionaryAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *dictionaryAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	return self.wordsDawg.parse(self, word, wordLower, seenParses)
}

func (self *dictionaryAnalyzer) getLexeme(form *Parse) []*Parse {
	return self.dict.getLexeme(self, form)
}

func (self *dictionaryAnalyzer) normalized(form *Parse) *Parse {
	return self.dict.normalized(self, form)
}

// метод tag "наследуется" от wordsDawg
