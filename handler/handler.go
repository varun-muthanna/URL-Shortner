package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Urls struct {
	url map[string]string
}

func NewURLShort() *Urls {
	return &Urls{url: make(map[string]string)}
}

func Generate(s string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var res string
	var val int

	//conversion to base 62 -> possible outcomes
	for i := 0; i < (len(s)); i++ {
		val = val + int(s[i])
	}

	for val > 0 {
		res = res + string(charset[val%62])
		val = val / 62
	}

	return res
}

func (u *Urls) ServeForm(w http.ResponseWriter, r *http.Request) {
	form := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>URL Shortner</title>
	</head>
	<body>
		<h1>Enter URL</h1>
		<form action="/shorten" method="post">
			<input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
		</form>
	</body>
	</html>
	`
	t := template.New("form")
	t, _ = t.Parse(form)
	t.Execute(w, nil)
}

func (u *Urls) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	longurl := r.FormValue("url")
	fmt.Println("Recieved %s", longurl)

	shorturl := Generate(longurl)
	u.url[shorturl] = longurl

	shortenedURL := fmt.Sprintf("http://localhost:3031/short/%s", shorturl)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Shortened URL: <a href="%s">%s</a></p>
	`, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)

}

func (u *Urls) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET SHORT CALLED")

	vars := mux.Vars(r) //stores as map,  of name of path --> key
	shortcode := vars["shortcode"]

	org, err := u.url[shortcode]

	if !err {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, org, http.StatusSeeOther)
}
