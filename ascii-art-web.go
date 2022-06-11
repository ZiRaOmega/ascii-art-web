package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Valeur string
}

func Get_ascii_char(caractere string, Font string) []string {
	output := []string{}
	for _, element := range caractere {
		char := rune(element - 32)
		file, err := os.Open(Font + ".txt")
		if err != nil {
			os.Exit(1)
			fmt.Println(err)
		}
		scanner := bufio.NewScanner(file)
		ascii := []string{}
		for scanner.Scan() {
			ascii = append(ascii, scanner.Text())
		}
		for i := char * 9; i < (char*9)+9; i++ {
			output = append(output, ascii[i])
		}
	}
	return output
}
func Show_ascii(ascii_char []string) string {
	nbr_lettre := len(ascii_char) / 9
	var word_array []string
	var result string
	for y := 1; y < 8; y++ {
		for i := y; i < nbr_lettre*9; i += 9 {
			word_array = append(word_array, ascii_char[i])
			if ascii_char[i] != "\n" {
				//fmt.Print(ascii_char[i])
				result += ascii_char[i]
			}
		}
		//fmt.Println("")
		result += "\n"
	}
	return result
}
func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/static/style.css" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		fmt.Println("Status code : 404")
		return
	}

	switch r.Method {
	case "GET":
		//
		http.ServeFile(w, r, "index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		inputText := r.FormValue("inputText")
		Font := r.FormValue("Font")
		to_print_slice := strings.Split(inputText, "\\n")
		for i := 0; i < len(to_print_slice); i++ {
			fmt.Fprintf(w, "%s", Show_ascii(Get_ascii_char(to_print_slice[i], Font)))
		}
		//fmt.Fprintf(w, "%s", Show_ascii(Get_ascii_char(inputText)))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
func style(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/style.css")
	return
}
func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
	return
}
func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/static/style.css", style)
	http.HandleFunc("/favicon.ico", favicon)
	fmt.Printf("Starting server for testing HTTP POST on http://localhost:8080 ...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
