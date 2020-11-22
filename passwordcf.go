package passwordcf

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func randomWordIndexes(num int, r *rand.Rand) []int {
	lines := make([]int, num)

	for i := 0; i < num; i++ {
		lines[i] = r.Intn(NumWords)
	}

	return lines
}

// returns a list of random words from wordlist
func passWords(numWords int) []string {

	// seed random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// word storage
	passWords := make([]string, numWords)

	// get a sorted list of indexes
	ridxs := randomWordIndexes(numWords, r)
	sort.Ints(ridxs)

	// map indexes to words
	for i := 0; i < numWords; i++ {
		passWords[i] = strings.Title(Words[ridxs[i]])
	}

	// shuffle order to make less predictable
	r.Shuffle(len(passWords), func(i, j int) { passWords[i], passWords[j] = passWords[j], passWords[i] })

	return passWords
}

// GeneratePassword API-entrypoint for generating passwords
func GeneratePassword(w http.ResponseWriter, r *http.Request) {
	// CORS preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Vary", "Origin")
		w.WriteHeader(http.StatusNoContent)
		return
	}	


	// make sure array of words is populated
	if Words == nil {
		initWords()
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// get requested number of words
	reqWords, err := strconv.Atoi(r.URL.Query().Get("numWords"))
	if err != nil || reqWords < 3 {
		//ignore invalid requests or too few words
		reqWords = 3
	}

	if reqWords > 10 {
		http.Error(w, "{ \"error\": \"numWords must be at least 3 and the maximum is 10\" }", 400)
		return
	}

	pwds := passWords(reqWords)

	fmt.Fprintf(w, "{ \"password\": \"%s\" }", strings.Join(pwds, " "))
}
