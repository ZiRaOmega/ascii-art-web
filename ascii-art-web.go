package main

import (
	"bufio"
	"color/color"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Page struct {
	Ascii       string
	Textareacol int
	Textarealin int
	Text        string
	Title       string
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
	var result string

	//Récupération Du Texte à Transformer et du Font choisi envoyé par la methode Post
	inputText := r.FormValue("inputText")
	Font := r.FormValue("Font")
	//fmt.Println(inputText)

	//Suppression des caractère de Retour a la ligne
	inputText = strings.ReplaceAll(inputText, `\r\n`, `\n`)

	//fmt.Println(inputText)
	inputText = strings.ReplaceAll(inputText, `\n`, string(rune(10)))
	to_print_slice := strings.Split(inputText, string(rune(10)))
	to_print_slice = append(to_print_slice, "\r")
	//fmt.Println(inputText)

	//fmt.Println(to_print_slice)
	// Création d'une nouvelle instance de template
	t := template.New("index")

	// Déclaration des fichiers à parser
	t = template.Must(t.ParseFiles("./templates/index.html", "./static/style.css", "./favicon.ico"))

	//Log console Status Code 400 Rouge, 200 Vert
	code := Status_code(w, r)
	if code == 400 || code == 404 || code == 403 {
		http.Error(w, strconv.Itoa(code)+" Bad Requests.", code)
		fmt.Println(color.ANSI_COLOR("RED") + strconv.Itoa(code) + " on " + r.RemoteAddr + r.URL.Path + color.ANSI_COLOR("RESET"))
		return
	} else {
		fmt.Println(color.ANSI_COLOR("GREEN") + strconv.Itoa(code) + " on " + r.RemoteAddr + r.URL.Path + color.ANSI_COLOR("RESET"))
	}

	//Mise en forme du Resultat en Ascii
	for i := 0; i < len(to_print_slice); i++ {
		if (i < len(to_print_slice[i])) && to_print_slice[i][(len(to_print_slice[i])-1):] != "\r" {
			result += Show_ascii(Get_ascii_char(to_print_slice[i], Font))
		} else if i < len(to_print_slice[i]) && to_print_slice[i][(len(to_print_slice[i])-1):] == "\r" {
			//to_print_slice[i] = strings.ReplaceAll(to_print_slice[i], string(rune(10)), "_")
			//to_print_slice[i] = to_print_slice[i][:len(to_print_slice[i])-1]
			result += Show_ascii(Get_ascii_char(to_print_slice[i][:len(to_print_slice[i])-1], Font))
		} else {
			re := regexp.MustCompile(`\r?\n`)
			to_print_slice[i] = re.ReplaceAllString(to_print_slice[i], ``)
			to_print_slice[i] = strings.ReplaceAll(to_print_slice[i], "\r", "")
			inputText = strings.ReplaceAll(inputText, "\r", "")
			result += Show_ascii(Get_ascii_char(to_print_slice[i], Font))
			//fmt.Println("error" + to_print_slice[i] + "l")
		}
	}

	//Creation d'une Page avec la valeur ascii et la taille du textarea avec un taille supérieur à 175 à l'aide de la struc Page
	col := (len(strings.Join(to_print_slice, "")) * 9)
	//fmt.Println(col)
	if col < 195 {
		col = 195
	}
	p := Page{result, col, len(to_print_slice) * 9, inputText, inputText}

	//On lance la template index avec la valeur P en valeur
	// La Page p sera réprésentée par le "." suivi de la variable Ex: p.Ascii et {{.Ascii}} dans notre Template
	// Exemple {{.}} == p
	err := t.ExecuteTemplate(w, "index", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}
func Status_code(w http.ResponseWriter, r *http.Request) int {

	if r.URL.Path == "/" || r.URL.Path == "/ascii-art" {
		return http.StatusOK
	} else if _, err := os.Stat("." + r.URL.Path); errors.Is(err, os.ErrNotExist) {
		return http.StatusNotFound
	} else {
		return http.StatusBadRequest
	}
}

func style(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/style.css")
	code := Status_code(w, r)
	if code == 400 || code == 404 || code == 403 {
		http.Error(w, "400 Bad Requests.", http.StatusBadRequest)
		fmt.Println(color.ANSI_COLOR("RED") + strconv.Itoa(code) + " on " + r.RemoteAddr + r.URL.Path + color.ANSI_COLOR("RESET"))
		return
	} else {
		fmt.Println(color.ANSI_COLOR("GREEN") + strconv.Itoa(code) + " on " + r.RemoteAddr + r.URL.Path + color.ANSI_COLOR("RESET"))
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
	fmt.Printf("Starting server for testing HTTP POST on http://localhost:8080 ...\n")

	//Commence le serveur sur le Port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
