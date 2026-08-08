package gomorphy

import (
	"fmt"
	"strings"
	"sync"
)

var (
	tagCls       *tagClass
	tagClassOnce sync.Once
)

func getTagClassInstance() *tagClass {
	tagClassOnce.Do(func() {
		tagCls = newTagClass()
	})
	return tagCls
}

type tagClass struct {
	partsOfSpeech             map[string]bool
	cases                     map[string]bool
	numbers                   map[string]bool
	genders                   map[string]bool
	animacy                   map[string]bool
	rareCases                 map[string]string
	knownGrammemes            map[string]string
	grammemeIncompatible      map[string][]string
	nonProductiveGrammemes    map[string]bool
	numeralAgreementGrammemes [5][]string
}

func newTagClass() *tagClass {

	tgCls := new(tagClass)

	tgCls.partsOfSpeech = map[string]bool{
		"NOUN": true, //  # имя существительное
		"ADJF": true, //  # имя прилагательное (полное)
		"ADJS": true, //  # имя прилагательное (краткое)
		"COMP": true, //  # компаратив
		"VERB": true, //  # глагол (личная форма)
		"INFN": true, //  # глагол (инфинитив)
		"PRTF": true, //  # причастие (полное)
		"PRTS": true, //  # причастие (краткое)
		"GRND": true, //  # деепричастие
		"NUMR": true, //  # числительное
		"ADVB": true, //  # наречие
		"NPRO": true, //  # местоимение-существительное
		"PRED": true, //  # предикатив
		"PREP": true, //  # предлог
		"CONJ": true, //  # союз
		"PRCL": true, //  # частица
		"INTJ": true, //  # междометие
	}

	tgCls.cases = map[string]bool{
		"nomn": true, //  # именительный падеж
		"gent": true, //  # родительный падеж
		"datv": true, //  # дательный падеж
		"accs": true, //  # винительный падеж
		"ablt": true, //  # творительный падеж
		"loct": true, //  # предложный падеж
		"voct": true, //  # звательный падеж
		"gen1": true, //  # первый родительный падеж
		"gen2": true, //  # второй родительный (частичный) падеж
		"acc2": true, //  # второй винительный падеж
		"loc1": true, //  # первый предложный падеж
		"loc2": true, //  # второй предложный (местный) падеж
	}

	tgCls.numbers = map[string]bool{
		"sing": true, //  # единственное число
		"plur": true, //  # множественное число
	}

	tgCls.genders = map[string]bool{
		"masc": true, //  # мужской род
		"femn": true, //  # женский род
		"neut": true, //  # средний род
	}

	tgCls.animacy = map[string]bool{
		"anim": true, //  # одушевлённое
		"inan": true, //  # неодушевлённое
	}

	tgCls.rareCases = map[string]string{
		"gen1": "gent",
		"gen2": "gent",
		"acc1": "accs",
		"acc2": "accs",
		"loc1": "loct",
		"loc2": "loct",
		"voct": "nomn",
	}

	// это видимо из файла grammemes.json
	tgCls.knownGrammemes = map[string]string{
		"POST": "ЧР", "NOUN": "СУЩ", "ADJF": "ПРИЛ", "ADJS": "КР_ПРИЛ", "COMP": "КОМП", "VERB": "ГЛ",
		"INFN": "ИНФ", "PRTF": "ПРИЧ", "PRTS": "КР_ПРИЧ", "GRND": "ДЕЕПР", "NUMR": "ЧИСЛ", "ADVB": "Н",
		"NPRO": "МС", "PRED": "ПРЕДК", "PREP": "ПР", "CONJ": "СОЮЗ", "PRCL": "ЧАСТ", "INTJ": "МЕЖД",
		"ANim": "Од-неод", "anim": "од", "inan": "неод", "GNdr": "хр", "masc": "мр", "femn": "жр",
		"neut": "ср", "ms-f": "мж", "NMbr": "Число", "sing": "ед", "plur": "мн", "Sgtm": "sg",
		"Pltm": "pl", "Fixd": "0", "CAse": "Падеж", "nomn": "им", "gent": "рд", "datv": "дт",
		"accs": "вн", "ablt": "тв", "loct": "пр", "voct": "зв", "gen1": "рд1", "gen2": "рд2",
		"acc2": "вн2", "loc1": "пр1", "loc2": "пр2", "Abbr": "аббр", "Name": "имя", "Surn": "фам",
		"Patr": "отч", "Geox": "гео", "Orgn": "орг", "Trad": "tm", "Subx": "субст?", "Supr": "превосх",
		"Qual": "кач", "Apro": "мест-п", "Anum": "числ-п", "Poss": "притяж", "V-ey": "*ею",
		"V-oy": "*ою", "Cmp2": "сравн2", "V-ej": "*ей", "ASpc": "Вид", "perf": "сов", "impf": "несов",
		"TRns": "Перех", "tran": "перех", "intr": "неперех", "Impe": "безл", "Impx": "безл?",
		"Mult": "мног", "Refl": "возвр", "PErs": "Лицо", "1per": "1л", "2per": "2л", "3per": "3л",
		"TEns": "Время", "pres": "наст", "past": "прош", "futr": "буд", "MOod": "Накл", "indc": "изъяв",
		"impr": "повел", "INvl": "Совм", "incl": "вкл", "excl": "выкл", "VOic": "Залог", "actv": "действ",
		"pssv": "страд", "Infr": "разг", "Slng": "жарг", "Arch": "арх", "Litr": "лит", "Erro": "опеч",
		"Dist": "искаж", "Ques": "вопр", "Dmns": "указ", "Prnt": "вводн", "V-be": "*ье", "V-en": "*енен",
		"V-ie": "*ие", "V-bi": "*ьи", "Fimp": "*несов", "Prdx": "предк?", "Coun": "счетн", "Coll": "собир",
		"V-sh": "*ши", "Af-p": "*предл", "Inmx": "не/одуш?", "Vpre": "в_предл", "Anph": "анаф",
		"Init": "иниц", "Adjx": "прил?", "Ms-f": "ор", "Hypo": "гипот", "NUMB": "ЧИСЛО", "intg": "цел",
		"real": "вещ", "PNCT": "ЗПР", "ROMN": "РИМ", "LATN": "ЛАТ", "UNKN": "НЕИЗВ"}

	tgCls.grammemeIncompatible = map[string][]string{
		"NOUN": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "PRTS"},
		"ADJF": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"ADJS": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "PRED", "NOUN", "PRTS"},
		"COMP": {"NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"VERB": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"INFN": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"PRTF": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"PRTS": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN"},
		"GRND": {"COMP", "NPRO", "ADVB", "PRCL", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"NUMR": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"ADVB": {"COMP", "NPRO", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"NPRO": {"COMP", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"PRED": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "NOUN", "PRTS"},
		"PREP": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "ADJS", "PRED", "NOUN", "PRTS"},
		"CONJ": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"PRCL": {"COMP", "NPRO", "ADVB", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "INTJ", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"INTJ": {"COMP", "NPRO", "ADVB", "PRCL", "GRND", "PRTF", "CONJ", "NUMR", "VERB", "INFN", "ADJF", "PREP", "ADJS", "PRED", "NOUN", "PRTS"},
		"anim": {"inan"},
		"inan": {"anim"},
		"masc": {"femn"},
		"femn": {"masc"},
		"neut": {"masc", "ms-f", "femn"},
		"ms-f": {"masc", "femn", "neut"},
		"sing": {"plur"},
		"plur": {"ms-f", "femn", "sing", "GNdr", "masc", "neut"},
		"nomn": {"accs", "gen2", "datv", "gent", "gen1", "ablt", "acc2", "loc1", "loct", "loc2", "voct"},
		"gent": {"accs", "gen2", "datv", "nomn", "gen1", "ablt", "acc2", "loc1", "loct", "loc2", "voct"},
		"datv": {"accs", "gen2", "gent", "nomn", "gen1", "ablt", "acc2", "loc1", "loct", "loc2", "voct"},
		"accs": {"gen2", "datv", "gent", "nomn", "gen1", "ablt", "acc2", "loc1", "loct", "loc2", "voct"},
		"ablt": {"accs", "gen2", "datv", "gent", "nomn", "gen1", "acc2", "loc1", "loct", "loc2", "voct"},
		"loct": {"accs", "gen2", "datv", "gent", "nomn", "gen1", "ablt", "acc2", "loc1", "loc2", "voct"},
		"gen1": {"gen2"},
		"gen2": {"gen1"},
		"loc1": {"loc2"},
		"loc2": {"loc1"},
		"perf": {"impf"},
		"impf": {"perf"},
		"tran": {"intr"},
		"intr": {"tran"},
		"1per": {"2per", "3per"},
		"2per": {"3per", "1per"},
		"3per": {"2per", "1per"},
		"pres": {"past", "futr"},
		"past": {"pres", "futr"},
		"futr": {"past", "pres"},
		"indc": {"impr"},
		"impr": {"indc"},
		"incl": {"excl"},
		"excl": {"incl"},
		"actv": {"pssv"},
		"pssv": {"actv"},
	}

	tgCls.nonProductiveGrammemes = map[string]bool{
		"NUMR": true, "NPRO": true, "PRED": true, "PREP": true,
		"CONJ": true, "PRCL": true, "INTJ": true, "Apro": true,
	}

	tgCls.numeralAgreementGrammemes = [5][]string{
		{"sing", "nomn"},
		{"sing", "accs"},
		{"sing", "gent"},
		{"plur", "nomn"},
		{"plur", "gent"},
	}

	return tgCls
}

// Replace rare cases (loc2/voct/...) with common ones (loct/nomn/...).
func (cls *tagClass) fixRareCases(grammemes map[string]bool) map[string]bool {
	newGrammes := make(map[string]bool, len(grammemes))
	for g := range grammemes {
		if v, ok := cls.rareCases[g]; ok {
			newGrammes[v] = true
		} else {
			newGrammes[g] = true
		}
	}
	return newGrammes
}

// Return Cyrillic representation for “tag_or_grammeme“ string
func (cls *tagClass) lat2cyr(tagOrGrammeme string) (string, error) {
	if tagOrGrammeme == "" {
		return "", nil
	}

	var tagsStrB strings.Builder
	tagsStrB.Grow(len(tagOrGrammeme))
	for part := range strings.SplitSeq(tagOrGrammeme, " ") {
		for tag := range strings.SplitSeq(part, ",") {
			if g, ok := cls.knownGrammemes[tag]; !ok {
				return "", fmt.Errorf("Unknown tag or grammeme: %#v", tag)
			} else {
				tagsStrB.WriteString(g)
				tagsStrB.WriteRune(',')
			}
		}
		tagsStrB.WriteRune(' ')
	}
	return strings.TrimRight(strings.ReplaceAll(tagsStrB.String(), ", ", " "), " "), nil
}

/*
dict_grammemes
(name, parent, alias, description)

['POST', '', 'ЧР', 'часть речи']
['NOUN', 'POST', 'СУЩ', 'имя существительное']
['ADJF', 'POST', 'ПРИЛ', 'имя прилагательное (полное)']
['ADJS', 'POST', 'КР_ПРИЛ', 'имя прилагательное (краткое)']
['COMP', 'POST', 'КОМП', 'компаратив']
['VERB', 'POST', 'ГЛ', 'глагол (личная форма)']
['INFN', 'POST', 'ИНФ', 'глагол (инфинитив)']
['PRTF', 'POST', 'ПРИЧ', 'причастие (полное)']
['PRTS', 'POST', 'КР_ПРИЧ', 'причастие (краткое)']
['GRND', 'POST', 'ДЕЕПР', 'деепричастие']
['NUMR', 'POST', 'ЧИСЛ', 'числительное']
['ADVB', 'POST', 'Н', 'наречие']
['NPRO', 'POST', 'МС', 'местоимение-существительное']
['PRED', 'POST', 'ПРЕДК', 'предикатив']
['PREP', 'POST', 'ПР', 'предлог']
['CONJ', 'POST', 'СОЮЗ', 'союз']
['PRCL', 'POST', 'ЧАСТ', 'частица']
['INTJ', 'POST', 'МЕЖД', 'междометие']
['ANim', '', 'Од-неод', 'категория одушевлённости']
['anim', 'ANim', 'од', 'одушевлённое']
['inan', 'ANim', 'неод', 'неодушевлённое']
['GNdr', '', 'хр', 'род / род не выражен']
['masc', 'ms-f', 'мр', 'мужской род']
['femn', 'ms-f', 'жр', 'женский род']
['neut', 'GNdr', 'ср', 'средний род']
['ms-f', 'GNdr', 'мж', 'общий род (м/ж)']
['NMbr', '', 'Число', 'число']
['sing', 'NMbr', 'ед', 'единственное число']
['plur', 'NMbr', 'мн', 'множественное число']
['Sgtm', '', 'sg', 'singularia tantum']
['Pltm', '', 'pl', 'pluralia tantum']
['Fixd', '', '0', 'неизменяемое']
['CAse', '', 'Падеж', 'категория падежа']
['nomn', 'CAse', 'им', 'именительный падеж']
['gent', 'CAse', 'рд', 'родительный падеж']
['datv', 'CAse', 'дт', 'дательный падеж']
['accs', 'CAse', 'вн', 'винительный падеж']
['ablt', 'CAse', 'тв', 'творительный падеж']
['loct', 'CAse', 'пр', 'предложный падеж']
['voct', 'nomn', 'зв', 'звательный падеж']
['gen1', 'gent', 'рд1', 'первый родительный падеж']
['gen2', 'gent', 'рд2', 'второй родительный (частичный) падеж']
['acc2', 'accs', 'вн2', 'второй винительный падеж']
['loc1', 'loct', 'пр1', 'первый предложный падеж']
['loc2', 'loct', 'пр2', 'второй предложный (местный) падеж']
['Abbr', '', 'аббр', 'аббревиатура']
['Name', '', 'имя', 'имя']
['Surn', '', 'фам', 'фамилия']
['Patr', '', 'отч', 'отчество']
['Geox', '', 'гео', 'топоним']
['Orgn', '', 'орг', 'организация']
['Trad', '', 'tm', 'торговая марка']
['Subx', '', 'субст?', 'возможна субстантивация']
['Supr', '', 'превосх', 'превосходная степень']
['Qual', '', 'кач', 'качественное']
['Apro', '', 'мест-п', 'местоименное']
['Anum', '', 'числ-п', 'порядковое']
['Poss', '', 'притяж', 'притяжательное']
['V-ey', '', '*ею', 'форма на -ею']
['V-oy', '', '*ою', 'форма на -ою']
['Cmp2', '', 'сравн2', 'сравнительная степень на по-']
['V-ej', '', '*ей', 'форма компаратива на -ей']
['ASpc', '', 'Вид', 'категория вида']
['perf', 'ASpc', 'сов', 'совершенный вид']
['impf', 'ASpc', 'несов', 'несовершенный вид']
['TRns', '', 'Перех', 'категория переходности']
['tran', 'TRns', 'перех', 'переходный']
['intr', 'TRns', 'неперех', 'непереходный']
['Impe', '', 'безл', 'безличный']
['Impx', '', 'безл?', 'возможно безличное употребление']
['Mult', '', 'мног', 'многократный']
['Refl', '', 'возвр', 'возвратный']
['PErs', '', 'Лицо', 'категория лица']
['1per', 'PErs', '1л', '1 лицо']
['2per', 'PErs', '2л', '2 лицо']
['3per', 'PErs', '3л', '3 лицо']
['TEns', '', 'Время', 'категория времени']
['pres', 'TEns', 'наст', 'настоящее время']
['past', 'TEns', 'прош', 'прошедшее время']
['futr', 'TEns', 'буд', 'будущее время']
['MOod', '', 'Накл', 'категория наклонения']
['indc', 'MOod', 'изъяв', 'изъявительное наклонение']
['impr', 'MOod', 'повел', 'повелительное наклонение']
['INvl', '', 'Совм', 'категория совместности']
['incl', 'INvl', 'вкл', 'говорящий включён (идем, идемте) ']
['excl', 'INvl', 'выкл', 'говорящий не включён в действие (иди, идите)']
['VOic', '', 'Залог', 'категория залога']
['actv', 'VOic', 'действ', 'действительный залог']
['pssv', 'VOic', 'страд', 'страдательный залог']
['Infr', '', 'разг', 'разговорное']
['Slng', '', 'жарг', 'жаргонное']
['Arch', '', 'арх', 'устаревшее']
['Litr', '', 'лит', 'литературный вариант']
['Erro', '', 'опеч', 'опечатка']
['Dist', '', 'искаж', 'искажение']
['Ques', '', 'вопр', 'вопросительное']
['Dmns', '', 'указ', 'указательное']
['Prnt', '', 'вводн', 'вводное слово']
['V-be', '', '*ье', 'форма на -ье']
['V-en', '', '*енен', 'форма на -енен']
['V-ie', '', '*ие', 'форма на -и- (веселие, твердостию); отчество с -ие']
['V-bi', '', '*ьи', 'форма на -ьи']
['Fimp', '', '*несов', 'деепричастие от глагола несовершенного вида']
['Prdx', '', 'предк?', 'может выступать в роли предикатива']
['Coun', '', 'счетн', 'счётная форма']
['Coll', '', 'собир', 'собирательное числительное']
['V-sh', '', '*ши', 'деепричастие на -ши']
['Af-p', '', '*предл', 'форма после предлога']
['Inmx', '', 'не/одуш?', 'может использоваться как одуш. / неодуш. ']
['Vpre', '', 'в_предл', 'Вариант предлога ( со, подо, ...)']
['Anph', '', 'анаф', 'Анафорическое (местоимение)']
['Init', '', 'иниц', 'Инициал']
['Adjx', '', 'прил?', 'может выступать в роли прилагательного']
['Ms-f', '', 'ор', 'колебание по роду (м/ж/с): кофе, вольво']
['Hypo', '', 'гипот', 'гипотетическая форма слова (победю, асфальтовее)']



cls._GRAMMEME_INCOMPATIBLE

POST frozenset()
NOUN frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'PRTS'})
ADJF frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
ADJS frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'PRED', 'NOUN', 'PRTS'})
COMP frozenset({'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
VERB frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
INFN frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
PRTF frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
PRTS frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN'})
GRND frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
NUMR frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
ADVB frozenset({'COMP', 'NPRO', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
NPRO frozenset({'COMP', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
PRED frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'NOUN', 'PRTS'})
PREP frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
CONJ frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
PRCL frozenset({'COMP', 'NPRO', 'ADVB', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'INTJ', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
INTJ frozenset({'COMP', 'NPRO', 'ADVB', 'PRCL', 'GRND', 'PRTF', 'CONJ', 'NUMR', 'VERB', 'INFN', 'ADJF', 'PREP', 'ADJS', 'PRED', 'NOUN', 'PRTS'})
ANim frozenset()
anim frozenset({'inan'})
inan frozenset({'anim'})
GNdr frozenset()
masc frozenset({'femn'})
femn frozenset({'masc'})
neut frozenset({'masc', 'ms-f', 'femn'})
ms-f frozenset({'masc', 'femn', 'neut'})
NMbr frozenset()
sing frozenset({'plur'})
plur frozenset({'ms-f', 'femn', 'sing', 'GNdr', 'masc', 'neut'})
Sgtm frozenset()
Pltm frozenset()
Fixd frozenset()
CAse frozenset()
nomn frozenset({'accs', 'gen2', 'datv', 'gent', 'gen1', 'ablt', 'acc2', 'loc1', 'loct', 'loc2', 'voct'})
gent frozenset({'accs', 'gen2', 'datv', 'nomn', 'gen1', 'ablt', 'acc2', 'loc1', 'loct', 'loc2', 'voct'})
datv frozenset({'accs', 'gen2', 'gent', 'nomn', 'gen1', 'ablt', 'acc2', 'loc1', 'loct', 'loc2', 'voct'})
accs frozenset({'gen2', 'datv', 'gent', 'nomn', 'gen1', 'ablt', 'acc2', 'loc1', 'loct', 'loc2', 'voct'})
ablt frozenset({'accs', 'gen2', 'datv', 'gent', 'nomn', 'gen1', 'acc2', 'loc1', 'loct', 'loc2', 'voct'})
loct frozenset({'accs', 'gen2', 'datv', 'gent', 'nomn', 'gen1', 'ablt', 'acc2', 'loc1', 'loc2', 'voct'})
voct frozenset()
gen1 frozenset({'gen2'})
gen2 frozenset({'gen1'})
acc2 frozenset()
loc1 frozenset({'loc2'})
loc2 frozenset({'loc1'})
Abbr frozenset()
Name frozenset()
Surn frozenset()
Patr frozenset()
Geox frozenset()
Orgn frozenset()
Trad frozenset()
Subx frozenset()
Supr frozenset()
Qual frozenset()
Apro frozenset()
Anum frozenset()
Poss frozenset()
V-ey frozenset()
V-oy frozenset()
Cmp2 frozenset()
V-ej frozenset()
ASpc frozenset()
perf frozenset({'impf'})
impf frozenset({'perf'})
TRns frozenset()
tran frozenset({'intr'})
intr frozenset({'tran'})
Impe frozenset()
Impx frozenset()
Mult frozenset()
Refl frozenset()
PErs frozenset()
1per frozenset({'2per', '3per'})
2per frozenset({'3per', '1per'})
3per frozenset({'2per', '1per'})
TEns frozenset()
pres frozenset({'past', 'futr'})
past frozenset({'pres', 'futr'})
futr frozenset({'past', 'pres'})
MOod frozenset()
indc frozenset({'impr'})
impr frozenset({'indc'})
INvl frozenset()
incl frozenset({'excl'})
excl frozenset({'incl'})
VOic frozenset()
actv frozenset({'pssv'})
pssv frozenset({'actv'})
Infr frozenset()
Slng frozenset()
Arch frozenset()
Litr frozenset()
Erro frozenset()
Dist frozenset()
Ques frozenset()
Dmns frozenset()
Prnt frozenset()
V-be frozenset()
V-en frozenset()
V-ie frozenset()
V-bi frozenset()
Fimp frozenset()
Prdx frozenset()
Coun frozenset()
Coll frozenset()
V-sh frozenset()
Af-p frozenset()
Inmx frozenset()
Vpre frozenset()
Anph frozenset()
Init frozenset()
Adjx frozenset()
Ms-f frozenset()
Hypo frozenset()
*/
