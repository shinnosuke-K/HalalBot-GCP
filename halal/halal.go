package halal

import (
	"log"
	"strings"
)

type halalFood struct {
	ngFoods []string
}

func New() *halalFood {
	return &halalFood{ngFoods: []string{"ワイン", "みりん", "日本酒", "ビール", "ラム酒", "料理酒", "豚肉", "豚", "ポーク", "ゼラチン", "ラード"}}
}

func (hf *halalFood) Judge(texts []string) string {
	for _, text := range texts {
		log.Println(text)
		if ok := hf.in(text); ok {
			return "impossible to eat"
		}
	}
	return "possible to eat"
}

func (hf *halalFood) in(word string) bool {
	for _, food := range hf.ngFoods {
		if ok := strings.Contains(word, food); ok {
			return true
		}
	}
	return false
}
