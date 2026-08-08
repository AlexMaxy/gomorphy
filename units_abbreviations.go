package gomorphy

import (
	"fmt"
	"strings"
)

// """
// Analyzer units for abbreviated words
// ------------------------------------
// """

type initialsAnalyzer struct {
	score      float32
	tagPattern string
	lettersSet map[string]bool
	selfTags   []string
}

func newInitialsAnalyzer(tagPattern string) *initialsAnalyzer {

	a := &initialsAnalyzer{
		score:      0.1,
		tagPattern: tagPattern,
		lettersSet: map[string]bool{},
	}

	for _, r := range "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЭЮЯ" {
		a.lettersSet[string(r)] = true
	}

	a.selfTags = a.getGenderCaseTags(a.tagPattern)

	return a
}

func (self *initialsAnalyzer) getGenderCaseTags(pattern string) []string {
	tags := make([]string, 0, 12)
	for _, gender := range [2]string{"masc", "femn"} {
		for _, case_ := range [6]string{"nomn", "gent", "datv", "accs", "ablt", "loct"} {
			tags = append(tags, fmt.Sprintf(pattern, gender, case_))
		}
	}
	return tags
}

func (self *initialsAnalyzer) parse(anal analyzer, word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	if !self.lettersSet[word] {
		return nil
	}

	parses := make([]*Parse, len(self.selfTags))
	for i, tag := range self.selfTags {
		parse := newParse(
			wordLower,
			tag,
			self.score,
		)
		parse.NormalForm = wordLower
		parse.MethodsStack = []Method{{Analyzer: anal, WordOrStack: word}}
		parses[i] = parse
	}
	return parses
}

func (self *initialsAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if !self.lettersSet[word] {
		return nil
	}

	tags := make([]*opencorporaTag, len(self.selfTags))
	for i, tag := range self.selfTags {
		tags[i] = newOpencorporaTag(tag)
	}
	return tags
}

type abbreviatedFirstNameAnalyzer struct {
	name     string
	terminal bool
	tagsMasc []string
	tagsFemn []string
	*initialsAnalyzer
}

func newAbbreviatedFirstNameAnalyzer() *abbreviatedFirstNameAnalyzer {
	tagPattern := "NOUN,anim,%s,Sgtm,Name,Fixd,Abbr,Init sing,%s"
	a := &abbreviatedFirstNameAnalyzer{
		name:             "AbbreviatedFirstNameAnalyzer",
		terminal:         false,
		initialsAnalyzer: newInitialsAnalyzer(tagPattern),
	}

	a.tagsMasc = make([]string, 0, len(a.selfTags))
	a.tagsFemn = make([]string, 0, len(a.selfTags))

	for _, tag := range a.selfTags {
		if strings.Contains(tag, "masc") {
			a.tagsMasc = append(a.tagsMasc, tag)
		} else if strings.Contains(tag, "femn") {
			a.tagsFemn = append(a.tagsFemn, tag)
		}
	}

	return a
}

func (self *abbreviatedFirstNameAnalyzer) Name() string {
	return self.name
}

func (self *abbreviatedFirstNameAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *abbreviatedFirstNameAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	return self.initialsAnalyzer.parse(self, word, wordLower, seenParses)
}

func (self *abbreviatedFirstNameAnalyzer) getLexeme(form *Parse) []*Parse {
	// # 2 lexemes: masc and femn
	var tags []string
	if form.Tag.Contains("masc") {
		tags = self.tagsMasc
	} else {
		tags = self.tagsFemn
	}

	parses := make([]*Parse, len(tags))
	for i, tag := range tags {
		parse := newParse(
			form.Word,
			tag,
			form.Score,
		)
		parse.NormalForm = form.NormalForm
		parse.MethodsStack = form.MethodsStack
		parses[i] = parse
	}

	return parses
}

func (self *abbreviatedFirstNameAnalyzer) normalized(form *Parse) *Parse {
	// # don't normalize female names to male names
	var tags []string
	if form.Tag.Contains("masc") {
		tags = self.tagsMasc
	} else {
		tags = self.tagsFemn
	}
	newForm := newParse(form.Word, tags[0], form.Score)
	newForm.NormalForm = form.NormalForm
	newForm.MethodsStack = form.MethodsStack
	return newForm
}

type abbreviatedPatronymicAnalyzer struct {
	name     string
	terminal bool
	*initialsAnalyzer
}

func newAbbreviatedPatronymicAnalyzer() *abbreviatedPatronymicAnalyzer {
	tagPattern := "NOUN,anim,%s,Sgtm,Patr,Fixd,Abbr,Init sing,%s"
	return &abbreviatedPatronymicAnalyzer{
		name:             "AbbreviatedPatronymicAnalyzer",
		terminal:         true,
		initialsAnalyzer: newInitialsAnalyzer(tagPattern),
	}
}

func (self *abbreviatedPatronymicAnalyzer) Name() string {
	return self.name
}

func (self *abbreviatedPatronymicAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *abbreviatedPatronymicAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	return self.initialsAnalyzer.parse(self, word, wordLower, seenParses)
}

func (self *abbreviatedPatronymicAnalyzer) getLexeme(form *Parse) []*Parse {
	parses := make([]*Parse, len(self.selfTags))
	for i, tag := range self.selfTags {
		parse := newParse(
			form.Word,
			tag,
			form.Score,
		)
		parse.NormalForm = form.NormalForm
		parse.MethodsStack = form.MethodsStack
		parses[i] = parse
	}

	return parses
}

// def normalized(self, form):
//     fixed_word, _, normal_form, score, methods_stack = form
//     return fixed_word, self._tags[0], normal_form, score, methods_stack

func (self *abbreviatedPatronymicAnalyzer) normalized(form *Parse) *Parse {
	newForm := newParse(form.Word, self.selfTags[0], form.Score)
	newForm.NormalForm = form.NormalForm
	newForm.MethodsStack = form.MethodsStack
	return newForm
}
