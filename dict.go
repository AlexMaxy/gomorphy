package gomorphy

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	precisionMask uint32 = 0xFFFF_FFFF
	isLeafBit     uint32 = 1 << 31
	hasLeafBit    uint32 = 1 << 8
	extensionBit  uint32 = 1 << 9
)

func unitHasLeaf(u uint32) bool { return u&hasLeafBit != 0 }
func unitLabel(u uint32) uint32 { return u & (isLeafBit | 0xFF) }
func unitOffset(u uint32) uint32 {
	return ((u >> 10) << ((u & extensionBit) >> 6)) & precisionMask
}
func unitValue(u uint32) uint32 { return u & (^isLeafBit & precisionMask) }

// dictionary is a read-only DAWG dictionary (array of 32-bit units).
type dictionary struct {
	units     []uint32
	paradigms [][]uint16 // paradigms[i] is a flat []uint16 of length N*3:
	//   [0:N]   -- suffix index for each form
	//   [N:2N]  -- gramtab tag ID for each form
	//   [2N:3N] -- paradigmPrefixes index for each form
	suffixes []string
	gramtab  []string // OpenCorpora tag string indexed by tag ID

	root uint32 // Root index

	minEndingFreq         int // значение из файла meta.json
	minParadigmPopularity int // значение из файла meta.json
	maxSuffixLength       int // значение из файла meta.json

	// # paradigm prefixes used for dictionary compilation
	// paradigmPrefixes are the three fixed paradigm prefixes used by pymorphy
	// Indices match meta.json → compile_options → paradigmPrefixes
	paradigmPrefixes [3]string
}

// read deserialises a dictionary from r (native-endian uint32 array).
func (d *dictionary) read(r io.Reader) error {
	var size uint32
	if err := binary.Read(r, binary.NativeEndian, &size); err != nil {
		return err
	}

	// TODO тут нужно извлечение значений из файла meta.json
	d.minEndingFreq = 2
	d.minParadigmPopularity = 3
	d.maxSuffixLength = 5
	d.paradigmPrefixes = [3]string{"", "по", "наи"}

	d.units = make([]uint32, size)
	return binary.Read(r, binary.NativeEndian, d.units)
}

func (d *dictionary) hasValue(index uint32) bool {
	return unitHasLeaf(d.units[index])
}

// followChar follows a single byte transition from index.
// Returns (next_index, true) or (0, false) if no such arc exists.
func (d *dictionary) followChar(label uint32, index uint32) (uint32, bool) {
	off := unitOffset(d.units[index])
	next := (index ^ off ^ label) & precisionMask
	if unitLabel(d.units[next]) != label {
		return 0, false
	}
	return next, true
}

// followBytes follows a sequence of bytes from index.
func (d *dictionary) followBytes(b []byte, index uint32) (uint32, bool) {
	for _, ch := range b {
		var ok bool
		index, ok = d.followChar(uint32(ch), index)
		if !ok {
			return 0, false
		}
	}
	return index, true
}

func (d *dictionary) followKey(b_key string, payloadSeparator []byte) (uint32, bool) {
	index, ok := d.followBytes([]byte(b_key), d.root)
	if !ok {
		return 0, false
	}
	index, ok = d.followBytes(payloadSeparator, index)
	if !ok {
		return 0, false
	}
	return index, true
}

// Build a normal form.
func (d *dictionary) buildNormalForm(f formParadigm, fixedWord string) string {
	if f.idx == 0 {
		return fixedWord
	}

	paradigm := d.paradigms[f.paraId]
	paradigmLen := len(paradigm) / 3

	stem := d.buildStem(paradigm, f.idx, fixedWord, d.suffixes)

	normalPrefixId := paradigm[paradigmLen*2+0]
	normalSuffixId := paradigm[0]

	normalPrefix := d.paradigmPrefixes[normalPrefixId]
	normalSuffix := d.suffixes[normalSuffixId]

	return normalPrefix + stem + normalSuffix
}

// Return word stem (given a word, paradigm and the word index).
func (d *dictionary) buildStem(paradigm []uint16, idx int, fixedWord string, suffixes []string) string {
	paradigmLen := len(paradigm) / 3

	prefixId := paradigm[paradigmLen*2+idx]
	prefix := d.paradigmPrefixes[prefixId]

	suffixId := paradigm[idx]
	suffix := suffixes[suffixId]

	if suffix != "" {
		return fixedWord[len(prefix) : len(fixedWord)-len(suffix)]
	}

	return fixedWord[len(prefix):]
}

// Return tag as a string.
func (d *dictionary) buildTagInfo(paraId int, idx int) string {
	paradigm := d.paradigms[paraId]
	tagInfoOffset := len(paradigm) / 3
	tagId := paradigm[tagInfoOffset+idx]
	return d.gramtab[tagId]
}

// This method assumes that DictionaryAnalyzer is the first and the only method in methods_stack.
func (d *dictionary) extractParaInfo(methodsStack []Method) (string, int, int) {
	analName := methodsStack[0].Analyzer.Name()
	if analName != "DictionaryAnalyzer" && analName != "FakeDictionary" {
		log.Fatal(fmt.Errorf(`[extractParaInfo] wrong methodsStack[0].Analyzer.Name(): %s`, analName))
	}
	return methodsStack[0].WordOrStack.(string), methodsStack[0].ParaIdOrStack.(int), methodsStack[0].Idx.(int)
}

func (d *dictionary) fixStack(anal analyzer, methodsStack []Method, word string, paraId int, idx int) []Method {
	method0 := Method{anal, word, paraId, idx}
	return append([]Method{method0}, methodsStack[1:]...)
}

// Return a lexeme (given a parsed word).
func (d *dictionary) getLexeme(anal analyzer, form *Parse) []*Parse {
	fixedWord := form.Word
	normalForm := form.NormalForm
	methodsStack := form.MethodsStack

	_, paraId, idx := d.extractParaInfo(methodsStack)

	para := d.paradigms[paraId]
	stem := d.buildStem(para, idx, fixedWord, d.suffixes)

	paradigm := d.buildParadigmInfo(paraId)

	result := make([]*Parse, len(paradigm))
	for index, pInfo := range paradigm {
		word := pInfo.prefix + stem + pInfo.suffix
		newMethodsStack := d.fixStack(anal, methodsStack, word, paraId, index)
		result[index] = &Parse{
			Word:         word,
			Tag:          newOpencorporaTag(pInfo.tag),
			NormalForm:   normalForm,
			Score:        1.0,
			MethodsStack: newMethodsStack,
		}
	}

	return result
}

func (d *dictionary) normalized(anal analyzer, form *Parse) *Parse {
	_, paraId, idx := d.extractParaInfo(form.MethodsStack)
	if idx == 0 {
		return form
	}
	tag := d.buildTagInfo(paraId, 0)
	newForm := newParse(form.NormalForm, tag, 1.0)
	newForm.NormalForm = form.NormalForm
	newForm.MethodsStack = d.fixStack(anal, form.MethodsStack, form.NormalForm, paraId, 0)

	return newForm
}

type paradigmInfo struct {
	prefix string
	tag    string
	suffix string
}

// Return a list of
//
//	(prefix, tag, suffix)
//
// tuples representing the paradigm.
func (d *dictionary) buildParadigmInfo(paraId int) []paradigmInfo {
	paradigm := d.paradigms[paraId]
	paradigmLen := len(paradigm) / 3

	res := make([]paradigmInfo, paradigmLen)
	for idx := range paradigmLen {
		prefixId := paradigm[paradigmLen*2+idx]
		prefix := d.paradigmPrefixes[prefixId]

		suffixId := paradigm[idx]
		suffix := d.suffixes[suffixId]

		res[idx] = paradigmInfo{
			prefix: prefix,
			tag:    d.buildTagInfo(paraId, idx),
			suffix: suffix,
		}
	}

	return res
}

// Exact matching (returns value)
func (d *dictionary) find(key []byte) int64 {
	index, ok := d.followBytes(key, d.root)
	if !ok {
		return -1
	}
	if !d.hasValue(index) {
		return -1
	}
	return d.value(index)
}

// Gets a value from a given index.
func (d *dictionary) value(index uint32) int64 {
	off := unitOffset(d.units[index])
	valueIndex := (index ^ off) & precisionMask
	return int64(unitValue(d.units[valueIndex])) // uint32 -> int64
}
