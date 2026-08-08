package gomorphy

import "fmt"

type seenParse struct {
	word   string
	tag    string
	paraId string // para_id м.б. или int или methods stack, поэтому тупо делаем строку из всего этого
}

func addParseIfNotSeen(parse *Parse, seenParses map[seenParse]bool) bool {

	reducedParse := seenParse{
		word: parse.Word,
		tag:  parse.Tag.RawTagsString,
	}

	if len(parse.MethodsStack) != 0 {
		reducedParse.paraId = fmt.Sprintf("%v", parse.MethodsStack[0].ParaIdOrStack)
	}

	if seenParses[reducedParse] {
		return false
	}

	seenParses[reducedParse] = true
	return true
}

func addTagIfNotSeen(tag *opencorporaTag, seenTags map[string]bool) bool {
	if seenTags[tag.RawTagsString] {
		return false
	}
	seenTags[tag.RawTagsString] = true
	return true
}

// Return a new form with “suffix“ attached
func withSuffix(form *Parse, suffix string) *Parse {
	parse := newParse(form.Word+suffix, form.Tag.RawTagsString, form.Score)
	parse.NormalForm = form.NormalForm + suffix
	parse.MethodsStack = form.MethodsStack
	return parse
}

//	Return a new form with “suffix_length“ chars removed from right
//
// suffix_length - in bytes
func withoutFixedSuffix(form *Parse, suffixLength int) *Parse {
	parse := newParse(
		form.Word[:len(form.Word)-suffixLength],
		form.Tag.RawTagsString,
		form.Score,
	)
	parse.NormalForm = form.NormalForm[:len(form.NormalForm)-suffixLength]
	parse.MethodsStack = form.MethodsStack

	return parse
}

//	Return a new form with “prefix_length“ chars removed from left
//
// prefix_length - in bytes
func withoutFixedPrefix(form *Parse, prefixLength int) *Parse {
	parse := newParse(form.Word[prefixLength:], form.Tag.RawTagsString, form.Score)
	parse.NormalForm = form.NormalForm[prefixLength:]
	parse.MethodsStack = form.MethodsStack
	return parse
}

// Return a new form with “prefix“ added
func withPrefix(form *Parse, prefix string) *Parse {
	parse := newParse(prefix+form.Word, form.Tag.RawTagsString, form.Score)
	parse.NormalForm = prefix + form.NormalForm
	parse.MethodsStack = form.MethodsStack
	return parse
}

// Return a new form with “methods_stack“ replaced with “new_methods_stack“
func replaceMethodsStack(form *Parse, newMethodsStack []Method) *Parse {
	parse := newParse(form.Word, form.Tag.RawTagsString, form.Score)
	parse.NormalForm = form.NormalForm
	parse.MethodsStack = newMethodsStack
	return parse
}

// Return a new form without last method from methods_stack
func withoutLastMethod(form *Parse) *Parse {
	parse := newParse(form.Word, form.Tag.RawTagsString, form.Score)
	parse.NormalForm = form.NormalForm
	parse.MethodsStack = form.MethodsStack[:len(form.MethodsStack)-1]
	return parse
}

// Return a new form with “method“ added to methods_stack
func appendMethod(form *Parse, method Method) *Parse {
	parse := newParse(form.Word, form.Tag.RawTagsString, form.Score)
	parse.NormalForm = form.NormalForm
	stack := form.MethodsStack
	stack = append(stack, method)
	parse.MethodsStack = stack
	return parse
}
