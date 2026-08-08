package gomorphy

type analogyAnalyzerUnit struct {
	// *baseAnalyzerUnit
}

func newAnalogyAnalyzerUnit() *analogyAnalyzerUnit {
	return &analogyAnalyzerUnit{
		// baseAnalyzerUnit: newBaseAnalyzerUnit(),
	}
}

func (self *analogyAnalyzerUnit) getLexeme(form *Parse) []*Parse {
	baseAnalyzer, thisMethod := self.methodInfo(form)
	f := withoutLastMethod(form)
	lexeme := baseAnalyzer.getLexeme(f)
	res := make([]*Parse, len(lexeme))
	for i, lex := range lexeme {
		res[i] = appendMethod(lex, thisMethod)
	}
	return res
}

func (self *analogyAnalyzerUnit) normalized(form *Parse) *Parse {
	baseAnalyzer, thisMethod := self.methodInfo(form)
	f := withoutLastMethod(form)
	normalForm := baseAnalyzer.normalized(f)
	return appendMethod(normalForm, thisMethod)
}

func (self *analogyAnalyzerUnit) methodInfo(form *Parse) (analyzer, Method) {
	methodsStack := form.MethodsStack
	methods := methodsStack[len(methodsStack)-2:]
	baseAnalyzer := methods[0].Analyzer
	return baseAnalyzer, methods[1]
}
