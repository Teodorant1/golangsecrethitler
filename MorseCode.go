package main

import (
	"encoding/json"
	"fmt"
	"github.com/gSpera/morse"
	"strings"
)

type MorseConverter struct {
}

func (MorseConverter MorseConverter) convertIntoMorseCode(text string) string {

	slicedText := strings.Split(text, " ")

	var slicedrealwords [][]string

	for _, currentword := range slicedText {

		currentword_symbols := strings.Split(currentword, "")
		slicedrealwords = append(slicedrealwords, currentword_symbols)

	}

	// we have a slice of slices, each subslice (one layer deep) represents a word , there are different pauses between letters and words
	//meaning we will have to return a slice of slices

	// this contains each word , we need each letter instead
	var MC_big_sentence [][][]string
	for _, subslice := range slicedrealwords {
		//	var SlicedMorseCodeWords []string
		var MC_words [][]string

		for _, letter := range subslice {
			fmt.Println(letter)
			morsecodeLetter := morse.ToMorse(letter)

			//array for 1 letter
			individualMorseSymbols := strings.Split(morsecodeLetter, "")

			MC_words = append(MC_words, individualMorseSymbols)
			//SlicedMorseCodeWords = append(SlicedMorseCodeWords, morsecodeLetter)
		}

		MC_big_sentence = append(MC_big_sentence, MC_words)

		//	MC_words = append(MC_words, SlicedMorseCodeWords)
	}

	marshal, err := json.Marshal(MC_big_sentence)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(marshal))

	//	morsecode := morse.ToMorse(text)
	//	//	fmt.Println(morsecode)
	//
	//	//	backtotext := morse.ToText(morsecode)
	//	//	fmt.Println(backtotext)
	//
	//	splitMorseSymbols := strings.Split(morsecode, "")
	//	fmt.Println(splitMorseSymbols)
	//
	//	fmt.Println("length of morse code", len(splitMorseSymbols))
	//
	//	for _, symbol := range splitMorseSymbols {
	//		fmt.Println("index is :", symbol)
	//	}
	//
	//	fmt.Println(splitMorseSymbols[1])
	//
	return string(marshal)
}
