package core

import "strconv"

// Contains function checks if an item exists in a list.
func Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

func HaveDifferent(list []string, list2 []string) bool {
	for _, v := range list {
		if !Contains(list2, v) {
			return true
		}
	}
	return false
}

func intListToStringList(intList []int) []string {
	stringList := make([]string, len(intList))
	for i, v := range intList {
		stringList[i] = strconv.Itoa(v)
	}
	return stringList
}
