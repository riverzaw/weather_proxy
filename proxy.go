package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	url, err := url.Parse("https://api.darksky.net/")
	if err != nil {
		panic(err)
	}
	username := os.Getenv("METEOLOGIN")
	password := os.Getenv("METEOPASS")
	key := os.Getenv("DARKSKYKEY")

	proxy := httputil.NewSingleHostReverseProxy(url)
	handler := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, pass, _ := r.BasicAuth()
			if user != username || pass != password {
				http.Error(w, "Wrong user or password", http.StatusUnauthorized)
				return
			}
			r.Host = url.Host
			r.URL.Path = "/forecast/" + key + r.URL.Path
			q := r.URL.Query()
			q.Add("exclude", "minutely,alerts,flags")
			q.Add("units", "si")
			r.URL.RawQuery = q.Encode()
			h.ServeHTTP(w, r)
		}
	}
	if err := http.ListenAndServe(":9034", handler(proxy)); err != nil {
		panic(err)
	}
}
