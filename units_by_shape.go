package gomorphy

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type shapeAnalyzer struct{}

func newShapeAnalyzer() *shapeAnalyzer {
	return &shapeAnalyzer{}
}

func (*shapeAnalyzer) getLexeme(form *Parse) []*Parse {
	return []*Parse{form}
}

func (self *shapeAnalyzer) normalized(form *Parse) *Parse {
	return form
}

// This analyzer marks integer numbers with "NUMB,int" or "NUMB,real" tags.
// Example: "12" -> NUMB,int; "12.4" -> NUMB,real
//
//	Don't confuse it with "NUMR": "тридцать" -> NUMR
type numberAnalyzer struct {
	name     string
	selfTags map[string]string
	terminal bool
	*shapeAnalyzer
}

func newNumberAnalyzer() *numberAnalyzer {
	return &numberAnalyzer{
		name:          "NumberAnalyzer",
		selfTags:      map[string]string{"intg": "NUMB,intg", "real": "NUMB,real"},
		terminal:      true,
		shapeAnalyzer: newShapeAnalyzer(),
	}
}

func (self *numberAnalyzer) Name() string {
	return self.name
}

func (self *numberAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *numberAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	wordTrimmed := strings.TrimSpace(word)
	if wordTrimmed == "" {
		return nil
	}

	isInt := false
	// if _, err := strconv.Atoi(wordTrimmed); err == nil {
	if self.isIntgNumber(wordTrimmed) {
		isInt = true
	} else if _, err := strconv.ParseFloat(strings.ReplaceAll(wordTrimmed, ",", "."), 64); err != nil {
		// strconv.ParseFloat в golang обрабатывает "1_2" как число.
		return nil
	}

	parse := &Parse{
		Word:         wordLower,
		NormalForm:   wordLower,
		Score:        0.9,
		MethodsStack: []Method{{Analyzer: self, WordOrStack: word}},
	}

	if isInt {
		parse.Tag = newOpencorporaTag(self.selfTags["intg"])

		return []*Parse{parse}
	}

	// Поведение как в питоне, "1_2_" или "_2_3" или 1__2 - UNKN,
	// но "1_2_3" - NUMB,intg
	// или "1_2.3_4" "1_2e3 - NUMB,real

	isFloat := false
	isLowbar := false
	for _, r := range wordTrimmed {
		if r == '.' || r == ',' || r == 'e' || r == 'E' {
			isFloat = true
			break
		} else if r == '_' {
			isLowbar = true
		}
	}

	if isLowbar && !isFloat {
		parse.Tag = newOpencorporaTag(self.selfTags["intg"])
	} else {
		parse.Tag = newOpencorporaTag(self.selfTags["real"])
	}

	return []*Parse{parse}
}

func (self *numberAnalyzer) isIntgNumber(s string) bool {
	for i, r := range s {
		if i == 0 && (r == '+' || r == '-') {
			if len(s) == 1 {
				return false
			}
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (self *numberAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if parses := self.parse(word, wordLower, nil); len(parses) != 0 {
		return []*opencorporaTag{parses[0].Tag}
	}
	return nil
}

// Return True if token looks like a Roman number:
//
//	>>> is_roman_number('II')
//	True
//	>>> is_roman_number('IX')
//	True
//	>>> is_roman_number('XIIIII')
//	False
//	>>> is_roman_number(”)
//	False
type romanNumberAnalyzer struct {
	name           string
	selfTag        string
	terminal       bool
	reRomanNumbers *regexp.Regexp
	*shapeAnalyzer
}

//     ^                   # start of string
//     M{0,4}              # thousands - 0 to 4 M's
//     (CM|CD|D?C{0,3})    # hundreds - 900 (CM), 400 (CD), 0-300 (0 to 3 C's),
//                         #            or 500-800 (D, followed by 0 to 3 C's)
//     (XC|XL|L?X{0,3})    # tens - 90 (XC), 40 (XL), 0-30 (0 to 3 X's),
//                         #        or 50-80 (L, followed by 0 to 3 X's)
//     (IX|IV|V?I{0,3})    # ones - 9 (IX), 4 (IV), 0-3 (0 to 3 I's),
//                         #        or 5-8 (V, followed by 0 to 3 I's)
//     $                   # end of string

func newRomanNumberAnalyzer() *romanNumberAnalyzer {
	return &romanNumberAnalyzer{
		name:           "RomanNumberAnalyzer",
		selfTag:        "ROMN",
		terminal:       false,
		reRomanNumbers: regexp.MustCompile(`(?i)^M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`),
		shapeAnalyzer:  newShapeAnalyzer(),
	}
}

func (self *romanNumberAnalyzer) Name() string {
	return self.name
}

func (self *romanNumberAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *romanNumberAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	if strings.TrimSpace(word) == "" {
		return nil
	} else if !self.reRomanNumbers.MatchString(word) {
		return nil
	}

	parse := &Parse{
		Word:         wordLower,
		Tag:          newOpencorporaTag(self.selfTag),
		NormalForm:   wordLower,
		Score:        0.9,
		MethodsStack: []Method{{Analyzer: self, WordOrStack: word}},
	}
	return []*Parse{parse}
}

func (self *romanNumberAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if parses := self.parse(word, wordLower, nil); len(parses) != 0 {
		return []*opencorporaTag{parses[0].Tag}
	}
	return nil
}

// This analyzer tags punctuation marks as "PNCT".
// Example: "," -> PNCT
// Return True if a word contains only spaces and punctuation marks
// and there is at least one punctuation mark:
//
//	>>> is_punctuation(', ')
//	True
//	>>> is_punctuation('..!')
//	True
//	>>> is_punctuation('x')
//	False
//	>>> is_punctuation(' ')
//	False
//	>>> is_punctuation('')
//	False
type punctuationAnalyzer struct {
	name     string
	selfTag  string
	terminal bool
	*shapeAnalyzer
}

func newPunctuationAnalyzer() *punctuationAnalyzer {
	return &punctuationAnalyzer{
		name:          "PunctuationAnalyzer",
		selfTag:       "PNCT",
		terminal:      true,
		shapeAnalyzer: newShapeAnalyzer(),
	}
}

func (self *punctuationAnalyzer) Name() string {
	return self.name
}

func (self *punctuationAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *punctuationAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	if strings.TrimSpace(word) == "" {
		return nil
	}
	for _, r := range word {
		if !unicode.IsSpace(r) && !unicode.Is(unicode.P, r) {
			return nil
		}
	}

	return []*Parse{&Parse{
		Word:         wordLower,
		Tag:          newOpencorporaTag(self.selfTag),
		NormalForm:   wordLower,
		Score:        0.9,
		MethodsStack: []Method{{Analyzer: self, WordOrStack: word}},
	}}
}

func (self *punctuationAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if parses := self.parse(word, wordLower, nil); len(parses) != 0 {
		return []*opencorporaTag{parses[0].Tag}
	}
	return nil
}

// This analyzer marks latin words with "LATN" tag.
// Example: "pdf" -> LATN
// Return True if all token letters are latin and there is at
// least one latin letter in the token:
//
//	>>> is_latin('foo')
//	True
//	>>> is_latin('123-FOO')
//	True
//	>>> is_latin('123')
//	False
//	>>> is_latin(':)')
//	False
//	>>> is_latin('')
//	False
type latinAnalyzer struct {
	name     string
	selfTag  string
	terminal bool
	*shapeAnalyzer
}

func newLatinAnalyzer() *latinAnalyzer {
	return &latinAnalyzer{
		name:          "LatinAnalyzer",
		selfTag:       "LATN",
		terminal:      true,
		shapeAnalyzer: newShapeAnalyzer(),
	}
}

func (self *latinAnalyzer) Name() string {
	return self.name
}

func (self *latinAnalyzer) isTerminal() bool {
	return self.terminal
}

func (self *latinAnalyzer) parse(word string, wordLower string, seenParses map[seenParse]bool) []*Parse {
	anyFlag := false
	allFlag := true
	for _, r := range word {
		if unicode.IsLetter(r) {
			anyFlag = true
			allFlag = allFlag && unicode.Is(unicode.Latin, r) && r != 'º' // 'º' - довесок для идентичности результатов
		}
	}

	if !(anyFlag && allFlag) {
		return nil
	}

	return []*Parse{&Parse{
		Word:         wordLower,
		Tag:          newOpencorporaTag(self.selfTag),
		NormalForm:   wordLower,
		Score:        0.9,
		MethodsStack: []Method{{Analyzer: self, WordOrStack: word}},
	}}
}

func (self *latinAnalyzer) tag(word, wordLower string, seenTags map[string]bool) []*opencorporaTag {
	if parses := self.parse(word, wordLower, nil); len(parses) != 0 {
		return []*opencorporaTag{parses[0].Tag}
	}
	return nil
}
