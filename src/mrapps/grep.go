package main

import "6.5840/mr"
import "strings"
import "fmt"
import "regexp"

const pattern string = "[mM]ap|[rR]educe"

func Map(filename string, contents string) []mr.KeyValue {
	kva := []mr.KeyValue{}
	rows := strings.Split(contents, "\n")
	for i, row := range rows {
		matched, _ := regexp.MatchString(pattern, row)
		if !matched {
			continue
		}
		v := fmt.Sprintf("%v:%v", filename, i+1)
		kva = append(kva, mr.KeyValue{row, v})
	}
	return kva
}

func Reduce(key string, values []string) string {
	return strings.Join(values, ", ")
}
