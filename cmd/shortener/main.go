package main

import (
	"io"
	"math/rand"
	"net/http"
)

var LocalHostValue string = "http://localhost:8080/"
var UrlDict = make(map[string]string)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	idSlice := []byte(id)
	idSlice = idSlice[1:]

	val, ok := UrlDict[string(idSlice)]
	if !ok {
		http.Error(w, "Original URL for this ID not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", val)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	ShortUrlResult := LocalHostValue + string(Shorter(b))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ShortUrlResult))
}

func Shorter(url []byte) []byte {
	letterBytes := string(url)
	b := make([]byte, 5)
	for i := range b {
		Elementvalue := letterBytes[rand.Intn(len(letterBytes))]
		if string(Elementvalue) != " " && string(Elementvalue) != "\n" && string(Elementvalue) != "/" {
			b[i] = Elementvalue
		} else {
			b[i] = '-'
		}
	}
	str := string(b)
	UrlDict[str] = string(url)
	return b
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetHandler(w, r)
	} else if r.Method == http.MethodPost {
		PostHandler(w, r)
	} else {
		http.Error(w, "Only GET/POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
