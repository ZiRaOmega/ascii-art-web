package main

import (
	"bufio"
	"color/color"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Page struct {
	Ascii       string
	Textareacol int
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
	for y := 1; y < 9; y++ {
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
func ascii_art(w http.ResponseWriter, r *http.Request) {
	Status_code(w, r)
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		inputText := r.FormValue("inputText")
		Font := r.FormValue("Font")
		to_print_slice := strings.Split(inputText, "\\n")
		for i := 0; i < len(to_print_slice); i++ {
			fmt.Fprintf(w, "%s", Show_ascii(Get_ascii_char(to_print_slice[i], Font)))
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// Création d'une page
	inputText := r.FormValue("inputText")
	Font := r.FormValue("Font")
	to_print_slice := strings.Split(inputText, "\\n")

	// Création d'une nouvelle instance de template
	t := template.New("index")

	// Déclaration des fichiers à parser
	t = template.Must(t.ParseFiles("./templates/index.html", "./static/style.css", "./favicon.ico"))

	// Exécution de la fusion et injection dans le flux de sortie
	// La variable p sera réprésentée par le "." dans le layout
	// Exemple {{.}} == p
	var result string
	code := Status_code(w, r)
	if code == 400 {
		http.Error(w, "400 Bad Requests.", http.StatusBadRequest)
		fmt.Println(color.ANSI_COLOR("RED") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
		return
	} else {
		fmt.Println(color.ANSI_COLOR("GREEN") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
	}
	for i := 0; i < len(to_print_slice); i++ {
		//fmt.Fprintf(w, "%s", Show_ascii(Get_ascii_char(to_print_slice[i], Font)))
		result += Show_ascii(Get_ascii_char(to_print_slice[i], Font))
	}
	p := Page{result, len(to_print_slice[0]) * 9}

	//fmt.Print(len(to_print_slice[0]))
	//fmt.Println(p.Ascii)
	err := t.ExecuteTemplate(w, "index", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}
func Status_code(w http.ResponseWriter, r *http.Request) int {
	if r.URL.Path == "/" || r.URL.Path == "/static/style.css" || r.URL.Path == "/favicon.ico" || r.URL.Path == "/ascii-art" {
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}

}

func style(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/style.css")
	code := Status_code(w, r)
	if code == 400 {
		http.Error(w, "400 Bad Requests.", http.StatusBadRequest)
		fmt.Println(color.ANSI_COLOR("RED") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
		return
	} else {
		fmt.Println(color.ANSI_COLOR("GREEN") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
	}
	return
}
func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
	return
}
func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/ascii-art", viewHandler)
	http.HandleFunc("/static/style.css", style)
	//http.HandleFunc("/favicon.ico", favicon)
	fmt.Printf("Starting server for testing HTTP POST on http://localhost:8080 ...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
