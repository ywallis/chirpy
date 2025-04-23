package main

import "strings"

var profanities = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func profanityFilter(input string, filter map[string]bool) string {

	output := []string{}
	for _, word := range strings.Split(input, " "){
		lcWord := strings.ToLower(word)
		if filter[lcWord] {
			output = append(output, "****")
		} else {
			output = append(output, word)
		} 
	}
	return strings.Join(output, " ")
}
