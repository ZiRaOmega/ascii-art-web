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
	Textarealin int
	Text        string
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
	for y := 0; y < 9; y++ {
		for i := y; i < nbr_lettre*9; i += 9 {
			word_array = append(word_array, ascii_char[i])
			if ascii_char[i] != "\n" {
				result += ascii_char[i]
			}
		}
		result += "\n"
	}
	return result
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//Récupération Du Texte à Transformer envoyé par la methode Post
	inputText := r.FormValue("inputText")

	Font := r.FormValue("Font")
	//Suppression des caractère de Retour a la ligne
	inputText = strings.Replace(inputText, `\n`, "\n", -1)
	inputText = strings.Replace(inputText, `\r`, "\n", -1)
	to_print_slice := strings.Split(inputText, "\n")
	// Création d'une nouvelle instance de template
	t := template.New("index")

	// Déclaration des fichiers à parser
	t = template.Must(t.ParseFiles("./templates/index.html", "./static/style.css", "./favicon.ico"))

	// Exécution de la fusion et injection dans le flux de sortie
	// La variable p sera réprésentée par le "." dans le layout
	// Exemple {{.}} == p
	var result string
	//Log console Status Code
	code := Status_code(w, r)
	if code == 400 {
		http.Error(w, "400 Bad Requests.", http.StatusBadRequest)
		fmt.Println(color.ANSI_COLOR("RED") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
		return
	} else {
		fmt.Println(color.ANSI_COLOR("GREEN") + strconv.Itoa(code) + " on " + "http://localhost:8080" + r.URL.Path + color.ANSI_COLOR("RESET"))
	}
	for i := 0; i < len(to_print_slice); i++ {
		if (i < len(to_print_slice[i])) && to_print_slice[i][(len(to_print_slice[i])-1):] != "\r" {
			result += Show_ascii(Get_ascii_char(to_print_slice[i], Font))
		} else if i < len(to_print_slice[i]) {
			to_print_slice[i] = to_print_slice[i][:len(to_print_slice[i])-1]
			result += Show_ascii(Get_ascii_char(to_print_slice[i], Font))

		}

	}
	//Creation d'une Page avec la valeur ascii et la taille du textarea
	col := len(strings.Join(to_print_slice, "")) * 5
	if len(strings.Join(to_print_slice, ""))*5 < 175 {
		col = 175
	}
	p := Page{result, col, len(to_print_slice) * 9, inputText}
	//On lance la template index avec la valeur P en valeur
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
	//HandleFunc Permet de definir les endpoints
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/ascii-art", viewHandler)
	http.HandleFunc("/static/style.css", style)
	//http.HandleFunc("/favicon.ico", favicon)
	fmt.Printf("Starting server for testing HTTP POST on http://localhost:8080 ...\n")
	//Commence le serveur sur le Port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
