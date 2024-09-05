package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Urls struct {
	redisclient *redis.Client
}

func NewURLShort() *Urls {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &Urls{redisclient: rdb}
}

func Generate(s string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var res string
	var val int

	//conversion to base 61 -> possible outcomes
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

	shorturl, err0 := u.redisclient.Get(context.Background(), longurl).Result()

	if err0 != redis.Nil {
		fmt.Printf("Already in redis instance")
	} else {
		shorturl = Generate(longurl)

		err1 := u.redisclient.Set(context.Background(), shorturl, longurl, 0).Err()
		err2 := u.redisclient.Set(context.Background(), longurl, shorturl, 0).Err()

		if err1 != nil {
			fmt.Printf("Could not set value in redis instance,%s", err.Error())
			return
		}

		if err2 != nil {
			fmt.Printf("Could not set value in redis instance,%s", err.Error())
			return
		}
	}

	shortenedURL := fmt.Sprintf("http://localhost:3031/short/%s", shorturl)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Shortened URL: <a href="%s">%s</a></p>
		<a href="/">Shorten Another URL</a> 
	`, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)

}

func (u *Urls) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET SHORT CALLED")

	vars := mux.Vars(r) //stores as map,  of name of path --> key
	shortcode := vars["shortcode"]

	org, err := u.redisclient.Get(context.Background(), shortcode).Result()

	if err != nil {
		http.Error(w, "URL not in redis", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, org, http.StatusSeeOther)
}
