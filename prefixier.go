package gomorphy

import (
	"sync"
)

type prefixWord struct {
	prefix         string
	unprefixedWord string
}

type suffixWord struct {
	unsuffixedWord string
	suffix         string
}

type prefixAnalyzer struct {
	*analogyAnalyzerUnit
}

func newPrefixAnalyzer() *prefixAnalyzer {
	return &prefixAnalyzer{
		analogyAnalyzerUnit: newAnalogyAnalyzerUnit(),
	}
}

func (self *prefixAnalyzer) getLexeme(form *Parse) []*Parse {
	baseAnalyzer, thisMethod := self.methodInfo(form)

	// TODO попробовать заменить на итератор/генератор
	// поискать по коду pymorphy все yield, они иногда выдают относительно большие списки.
	prefix := thisMethod.WordOrStack.(string)
	f := withoutFixedPrefix(form, len(prefix))
	f = withoutLastMethod(f)
	lexeme := baseAnalyzer.getLexeme(f)
	result := make([]*Parse, len(lexeme))
	for i, lex := range lexeme {
		result[i] = appendMethod(withPrefix(lex, prefix), thisMethod)
	}
	return result
}

func (self *prefixAnalyzer) normalized(form *Parse) *Parse {
	baseAnalyzer, thisMethod := self.methodInfo(form)
	prefix := thisMethod.WordOrStack.(string)
	f := withoutFixedPrefix(form, len(prefix))
	f = withoutLastMethod(f)
	normalForm := baseAnalyzer.normalized(f)
	return appendMethod(withPrefix(normalForm, prefix), thisMethod)
}

type prefixMatcher struct {
	*prefixTrie
}

var (
	pfxMatcher *prefixMatcher
	pfxOnce    sync.Once
)

func getPrefixMatcherInstance() *prefixMatcher {
	pfxOnce.Do(func() {
		pfxMatcher = newPrefixMatcher()
	})
	return pfxMatcher
}

func newPrefixMatcher() *prefixMatcher {
	prefixes := [...]string{
		"авиа", "авто", "аква", "анти", "анти-", "антропо", "архи", "арт", "арт-", "астро", "аудио", "аэро",
		"без", "бес", "био", "вело", "взаимо", "вне", "внутри", "видео", "вице-", "вперед", "впереди", "гекто",
		"гелио", "гео", "гетеро", "гига", "гигро", "гипер", "гипо", "гомо", "дву", "двух", "де", "дез",
		"дека", "деци", "дис", "до", "евро", "за", "зоо", "интер", "инфра", "квази", "квази-", "кило",
		"кино", "контр", "контр-", "космо", "космо-", "крипто", "лейб-", "лже", "лже-", "макро", "макси", "макси-",
		"мало", "меж", "медиа", "медиа-", "мега", "мета", "мета-", "метео", "метро", "микро", "милли", "мини",
		"мини-", "моно", "мото", "много", "мульти", "нано", "нарко", "не", "небез", "недо", "нейро", "нео",
		"низко", "обер-", "обще", "одно", "около", "орто", "палео", "пан", "пара", "пента", "пере", "пиро",
		"поли", "полу", "после", "пост", "пост-", "порно", "пра", "пра-", "пред", "пресс-", "противо", "противо-",
		"прото", "псевдо", "псевдо-", "радио", "разно", "ре", "ретро", "ретро-", "само", "санти", "сверх", "сверх-",
		"спец", "суб", "супер", "супер-", "супра", "теле", "тетра", "топ-", "транс", "транс-", "ультра", "унтер-",
		"штаб-", "экзо", "эко", "эндо", "эконом-", "экс", "экс-", "экстра", "экстра-", "электро", "энерго", "этно",
	}

	pm := &prefixMatcher{
		prefixTrie: newPrefixTrie(),
	}

	for _, pfx := range prefixes {
		pm.prefixTrie.insert(pfx)
	}

	return pm
}

// func (self *prefixMatcher) getPrefixes(word string) []string {
// 	prefixes := []string{}
// 	for _, prefix := range self.prefixes {
// 		if strings.HasPrefix(word, prefix) {
// 			prefixes = append(prefixes, prefix)
// 		}
// 	}
// 	return prefixes
// }

// func (self *prefixMatcher) isPrefixed(word string) bool {
// 	for _, prefix := range self.prefixes {
// 		if strings.HasPrefix(word, prefix) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (self *prefixMatcher) getPrefixes(word string) []string {
	return self.prefixTrie.findPrefixes(word)
}

func (self *prefixMatcher) isPrefixed(word string) bool {
	return len(self.prefixTrie.findPrefixes(word)) != 0
}

// префиксное дерево

type pfxNode struct {
	children map[rune]*pfxNode
	isEnd    bool
	prefix   string
}

type prefixTrie struct {
	root *pfxNode
}

func newPrefixTrie() *prefixTrie {
	return &prefixTrie{root: &pfxNode{children: make(map[rune]*pfxNode)}}
}

// insert добавляет префикс в дерево
func (t *prefixTrie) insert(prefix string) {
	current := t.root
	for _, char := range prefix {
		if _, exists := current.children[char]; !exists {
			current.children[char] = &pfxNode{children: make(map[rune]*pfxNode)}
		}
		current = current.children[char]
	}
	current.isEnd = true
	current.prefix = prefix
}

// findPrefixes находит все подходящие префиксы для слова
func (t *prefixTrie) findPrefixes(word string) []string {
	var matches []string
	current := t.root

	for _, char := range word {
		if nextNode, exists := current.children[char]; exists {
			current = nextNode
			if current.isEnd {
				matches = append(matches, current.prefix)
			}
		} else {
			break
		}
	}
	return matches
}

/*
из украинской части кода
{"2D-", "2G-", "3D-", "3G-", "4D-", "4G-", "CAD-", "call-", "CD-", "CDMA-", "CFI-", "CNG-", "DDoS-", "DNS-", "DoS-",
"DSL-", "dvd-", "e-", "fashion-", "FM-", "ftp-", "G-", "GMP-", "GPRS-", "GPS-", "grid-", "GSM-", "HD-", "HR-",
"HSDPA-", "ID-", "IMEA-", "IP-", "IT-", "led-", "LCD-", "LNG-", "live-", "MLM-", "MTV-", "mp3-", "n-", "OSB-",
"pdf-", "PhD-", "PIN-", "POS-", "pr-", "QR-", "R'n'B-", "R'N'B-", "R&B-", "R&D-", "s-", "sim-", "SOS-", "SPA-",
"sms-", "TV-", "UMTS-", "USB-", "VIN-", "vip-", "VoIP-", "WAP-", "web-", "X-", "Y-", "аль-", "альфа-", "анти-",
"АРВ-", "арт-", "аудіо-", "байк-", "байкер-", "бард-", "бас-", "бета-", "бізнес-", "бліц-", "блог-", "блок-",
"блюз-", "бомж-", "бонус-", "ботокс-", "боулінг-", "брейк-", "бренд-", "бундес-", "вакуум-", "веб-", "велнес-",
"ВІЛ-", "віп-", "віце-", "гала-", "гамма-", "гей-", "гейм-", "генерал-", "гештальт-", "ГМ-", "ГМО-", "гольф-",
"гоп-", "горе-", "готик-", "гранд-", "ґранд-", "графіті-", "грид-", "грумінг-", "дайв-", "дайвінг-", "данс-",
"даун-", "дельта-", "денс-", "дзен-", "джаз-", "диво-", "дизайн-", "дизель-", "долбі-", "допінг-", "ДОТС-", "драг-",
"дрес-", "дубль-", "дурман-", "е-", "екіпаж-", "економ-", "експерт-", "екс-", "експрес-", "екстра-", "екстрим-",
"екшн-", "еліт-", "ерзац-", "ескорт-", "євро-", "жлоб-", "зіц-", "зомбі-", "ЗПГ-", "івент-", "імідж-", "інвест-",
"інді-", "інсентив-", "інтернет-", "інтим-", "інформ-", "історико-", "ІТ-", "ІЧ-", "йога-", "камер-", "кантрі-",
"караоке-", "кастинг-", "квазі-", "кемпінг-", "кваліфайн-", "кібер-", "кітч-", "козак-", "коктейль-", "колл-",
"комік-", "комікс-", "майстер-", "конгрес-", "консалтинг-", "контент-", "контр-", "конференц-", "концепт-",
"кредит-", "кремль-", "крос-", "КСВ-", "лайт-", "лаунж-", "лейб-", "лесбі-", "лгбт-", "лже-", "ліберал-", "лор-",
"люкс-", "люмпен-", "максі-", "маркетинг-", "мас-", "мега-", "медіа-", "менеджмент-", "метал-", "міді-", "мікс-",
"мілітарі-", "міні-", "МММ-", "модерн-", "мульт-", "мультимедіа-", "напів-", "націонал-", "нація-", "НВЧ-",
"нокаут-", "ностальжі-", "нью-", "обер-", "онлайн-", "офіс-", "ОУН-", "панк-", "ПВХ-", "ПЕТ-", "піар-", "пін-",
"плейбек-", "ПЛР-", "покер-", "поп-", "пост-", "поттер-", "постпродакшн-", "прайм-", "прайс-", "прем'єр-",
"преміум-", "прес-", "приват-", "продакшн-", "профі-", "псевдо-", "реаліті-", "реггі-", "резус-", "рейв-",
"рентген-", "рейтинг-", "реп-", "ретро-", "референс-", "референц-", "ритм-", "РК-", "рок-", "ротарі-", "РХБ-",
"салон-", "саунд-", "своп-", "секонд-", "секс-", "сексі-", "сервіс-", "скейт-", "скінхед-", "скретч-", "слем-",
"смарт-", "смс-", "СНІД-", "соціал-", "СОС-", "соул-", "софт-", "спа-", "спам-", "спаринг-", "СПГ-", "спорт-",
"спрей-", "стартап-", "стоп-", "стрес-", "стрип-", "стриптиз-", "супер-", "тайм-", "талант-", "тандем-", "танц-",
"тату-", "ТБ-", "телеком-", "тест-", "топ-", "топлес-", "торент-", "тренд-", "тренінг-", "треш-", "триб'ют-",
"трофі-", "тур-", "тюнинг-", "УЗД-", "ура-", "УФ-", "фан-", "фест-", "фешн-", "фітнес-", "флеш-", "ФМ-", "фолк-",
"фольк-", "хеш-", "цар-", "чудо-", "хайтек-", "хард-", "хіпі-", "хостел-", "чіп-", "шейпінг-", "шенген-", "шеф-",
"шопінг-", "шоу-", "штаб-", "юніор-", "авіа", "авто", "аква", "анти", "антропо", "архі", "арт", "астро", "аудіо",
"аеро", "без", "біо", "вело", "взаємо", "поза", "внутрішньо", "відео", "вперед", "гекто", "гелио", "гео", "гетеро",
"гіга", "гігро", "гіпер", "гіпо", "гомо", "дво", "де", "дез", "деци", "дис", "до", "євро", "за", "зоо", "інтер",
"інфра", "квазі", "кіло", "кіно", "контр", "космо", "космо-", "крипто", "лже", "макро", "мало", "між", "медіа",
"мега", "мета", "мета-", "метео", "метро", "мікро", "мілі", "міні", "моно", "мото", "багато", "нано", "нарко", "не",
"нейро", "нео", "низько", "загально", "навколо", "орто", "палео", "пан", "пара", "пента", "пере", "піро", "полі",
"полу", "після", "пост", "порно", "пра", "пра-", "перед", "проти", "проти-", "прото", "псевдо", "радіо", "разно",
"ре", "ретро", "само", "санти", "над", "над-", "спец", "суб", "супер", "теле", "тетра", "транс", "транс-", "ультра",
"унтер-", "екзо", "еко", "ендо", "екс", "екстра", "електро", "енерго", "етно"}

*/
