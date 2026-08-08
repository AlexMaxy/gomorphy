package gomorphy

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	morphAnalyzer *MorphAnalyzer
	morphOnce     sync.Once
	morphErr      error
)

// Default returns the shared Analyzer loaded from embedded dictionary data
// The dictionary is initialised on the first call and cached; subsequent calls
// return the same instance. Safe for concurrent use
func GetMorphInstance(opencorporaPath string) (*MorphAnalyzer, error) {
	morphOnce.Do(func() {
		morphAnalyzer, morphErr = initMorphAnalyzer(opencorporaPath)
	})
	return morphAnalyzer, morphErr
}

type compiledReplaces struct {
	bReplaceChar []byte
	uReplaceChar string
}

type MorphAnalyzer struct {
	words           wordsDawg
	scores          intDawg
	predSuffixes    predictionSuffixesDawgs
	analyzers       []analyzer                // список анализаторов
	charSubstitutes map[rune]compiledReplaces // замена е на ё
	*tagClass
}

// НЕ использовать как самостоятельный конструктор.
// Там внутри завязки на глобальную переменную morphAnalyzer.
func initMorphAnalyzer(opencorporaPath string) (*MorphAnalyzer, error) {
	a := &MorphAnalyzer{
		analyzers: make([]analyzer, 0, 14),
		charSubstitutes: map[rune]compiledReplaces{'е': {
			bReplaceChar: []byte{209, 145},
			uReplaceChar: "ё",
		}},
		tagClass: getTagClassInstance(),
	}

	raw, err := os.ReadFile(filepath.Join(opencorporaPath, "words.dawg"))
	if err != nil {
		return nil, err
	}
	if err := a.words.load(bytes.NewReader(raw)); err != nil {
		return nil, err
	}

	raw, err = os.ReadFile(filepath.Join(opencorporaPath, "p_t_given_w.intdawg"))
	if err != nil {
		return nil, err
	}
	if err := a.scores.load(bytes.NewReader(raw)); err != nil {
		return nil, err
	}

	if a.predSuffixes, err = newPredictionSuffixesDawgs(opencorporaPath); err != nil {
		return nil, err
	}

	// paradigms.array: uint16 LE count, then per paradigm: uint16 LE length + data
	raw, err = os.ReadFile(filepath.Join(opencorporaPath, "paradigms.array"))
	if err != nil {
		return nil, err
	}
	if err := a.loadParadigms(raw); err != nil {
		return nil, err
	}

	raw, err = os.ReadFile(filepath.Join(opencorporaPath, "suffixes.json"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &a.words.dict.suffixes); err != nil {
		return nil, err
	}

	raw, err = os.ReadFile(filepath.Join(opencorporaPath, "gramtab-opencorpora-int.json"))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &a.words.dict.gramtab); err != nil {
		return nil, err
	}

	// _Unit(analyzer=DictionaryAnalyzer(), terminal=False)
	// _Unit(analyzer=AbbreviatedFirstNameAnalyzer(letters='АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЭЮЯ', score=0.1, tag_pattern='NOUN,anim,%(gender)s,Sgtm,Name,Fixd,Abbr,Init sing,%(case)s'), terminal=False)
	// _Unit(analyzer=AbbreviatedPatronymicAnalyzer(letters='АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЭЮЯ', score=0.1, tag_pattern='NOUN,anim,%(gender)s,Sgtm,Patr,Fixd,Abbr,Init sing,%(case)s'), terminal=True)
	// _Unit(analyzer=NumberAnalyzer(score=0.9), terminal=True)
	// _Unit(analyzer=PunctuationAnalyzer(score=0.9), terminal=True)
	// _Unit(analyzer=RomanNumberAnalyzer(score=0.9), terminal=False)
	// _Unit(analyzer=LatinAnalyzer(score=0.9), terminal=True)
	// _Unit(analyzer=HyphenSeparatedParticleAnalyzer(particles_after_hyphen=['-то', '-ка', '-таки', '-де', '-тко', '-тка', '-с', '-ста'], score_multiplier=0.9), terminal=True)
	// _Unit(analyzer=HyphenAdverbAnalyzer(score_multiplier=0.7), terminal=True)
	// _Unit(analyzer=HyphenatedWordsAnalyzer(score_multiplier=0.75, skip_prefixes=<...>), terminal=True)
	// _Unit(analyzer=KnownPrefixAnalyzer(known_prefixes=<...>, min_remainder_length=3, score_multiplier=0.75), terminal=True)
	// _Unit(analyzer=UnknownPrefixAnalyzer(score_multiplier=0.5), terminal=False)
	// _Unit(analyzer=KnownSuffixAnalyzer(min_word_length=4, score_multiplier=0.5), terminal=True)
	// _Unit(analyzer=UnknAnalyzer(), terminal=True)

	a.addAnalyzer(newDictionaryAnalyzer(&a.words))      // 0
	a.addAnalyzer(newAbbreviatedFirstNameAnalyzer())    // 1
	a.addAnalyzer(newAbbreviatedPatronymicAnalyzer())   // 2
	a.addAnalyzer(newNumberAnalyzer())                  // 3
	a.addAnalyzer(newPunctuationAnalyzer())             // 4
	a.addAnalyzer(newRomanNumberAnalyzer())             // 5
	a.addAnalyzer(newLatinAnalyzer())                   // 6
	a.addAnalyzer(newHyphenSeparatedParticleAnalyzer()) // 7
	a.addAnalyzer(newHyphenAdverbAnalyzer())            // 8
	a.addAnalyzer(newHyphenatedWordsAnalyzer())         // 9
	a.addAnalyzer(newKnownPrefixAnalyzer())             // 10
	a.addAnalyzer(newUnknownPrefixAnalyzer())           // 11
	a.addAnalyzer(newKnownSuffixAnalyzer(&a.words))     // 12
	a.addAnalyzer(newUnknAnalyzer())                    // 13

	return a, nil
}

func (m *MorphAnalyzer) loadParadigms(raw []byte) error {
	r := bytes.NewReader(raw)

	var n uint16
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return err
	}
	m.words.dict.paradigms = make([][]uint16, n)
	for i := range m.words.dict.paradigms {
		var length uint16
		if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
			return err
		}
		para := make([]uint16, length)
		if err := binary.Read(r, binary.LittleEndian, para); err != nil {
			return err
		}
		m.words.dict.paradigms[i] = para
	}
	return nil
}

// Lat2Cyr возвращает кириллическое представление тега или граммемы. Возвращает ошибку, если тег неизвестен.
func (m *MorphAnalyzer) Lat2Cyr(tagOrGrammeme string) (string, error) {
	return m.tagClass.lat2cyr(tagOrGrammeme)
}

// NormalForms возвращает список нормальных форм слова.
func (m *MorphAnalyzer) NormalForms(word string) []string {
	parses := m.Parse(word)
	seen := make(map[string]struct{}, len(parses))
	result := make([]string, 0, len(parses))
	for _, p := range parses {
		if _, ok := seen[p.NormalForm]; !ok {
			seen[p.NormalForm] = struct{}{}
			result = append(result, p.NormalForm)
		}
	}
	return result
}

// Check if a “word“ is in the dictionary.
//
// By default, some fuzziness is allowed, depending on a
// dictionary - e.g. for Russian ё letters replaced with е are handled.
// Pass “strict=True“ to make matching strict (e.g. if it is
// guaranteed the “word“ has correct е/ё or г/ґ letters).
//
// .. note::
//
//	Dictionary words are not always correct words;
//	the dictionary also contains incorrect forms which
//	are commonly used. So for spellchecking tasks this
//	method should be used with extra care.
func (m *MorphAnalyzer) WordIsKnown(word string, strict bool) bool {
	var substitutesCompiled map[rune]compiledReplaces
	if strict {
		// pass
	} else {
		substitutesCompiled = m.charSubstitutes
	}
	return m.words.wordIsKnown(strings.ToLower(word), substitutesCompiled)
}
