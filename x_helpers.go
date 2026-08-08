package gomorphy

import "unicode/utf8"

// isSubset (<=) проверяет, является ли «левое» множество подмножеством «правого».
func isSubset(subset map[string]bool, set map[string]bool) bool {
	if len(subset) > len(set) {
		return false
	}
	for s := range subset {
		if !set[s] {
			return false
		}
	}
	return true
}

// Пересечение множеств (&).
// Возвращают новую коллекцию, содержащую только те элементы, что есть в двух исходных множествах.
// Если таковых не найдётся, получится пустая коллекция.
func intersection(set1 map[string]bool, set2 map[string]bool) map[string]bool {
	intersected := map[string]bool{}
	for key := range set1 {
		if set2[key] {
			intersected[key] = true
		}
	}
	return intersected
}

// Тоже, что и intersection, только возвращает длину пересечения
func intersectionLen(set1 map[string]bool, set2 map[string]bool) int {
	n := 0
	for key := range set1 {
		if set2[key] {
			n++
		}
	}
	return n
}

// Симметрическая разность (^) — это набор элементов, которые принадлежат либо первому,
// либо второму множеству, но не их пересечению.
// Иными словами, симметрическая разность содержит все элементы обоих множеств, кроме общих.
func symmetricDifference(set1 map[string]bool, set2 map[string]bool) map[string]bool {
	diff := map[string]bool{}
	for key := range set1 {
		if !set2[key] {
			diff[key] = true
		}
	}
	for key := range set2 {
		if !set1[key] {
			diff[key] = true
		}
	}
	return diff
}

// тоже, что и symmetricDifference, только возвращает длину множества
func symmetricDifferenceLen(set1 map[string]bool, set2 map[string]bool) int {
	n := 0
	for key := range set1 {
		if !set2[key] {
			n++
		}
	}
	for key := range set2 {
		if !set1[key] {
			n++
		}
	}
	return n
}

// stringSliceByRuneIndex возвращает срез строки по индексам рун start/end, без выделения памяти.
//
// Если startRune = -1, то срез s[ : endRune ]
//
// Eсли endRune = -1, то срез s[ startRune : ]
func stringSliceByRuneIndex(s string, startRune, endRune int) string {
	if startRune >= utf8.RuneCountInString(s) {
		return ""
	}

	startByte, endByte := 0, len(s)
	runeCount := 0

	for i := range s {
		if runeCount == startRune {
			startByte = i
		}
		if runeCount == endRune {
			endByte = i
			break
		}
		runeCount++
	}

	return s[startByte:endByte]
}

// runeAt возвращает руну по индексу рун
func runeAt(s string, runeIdx int) (rune, bool) {
	currentRuneIdx := 0
	for _, r := range s {
		if currentRuneIdx == runeIdx {
			return r, true
		}
		currentRuneIdx++
	}
	return 0, false
}
