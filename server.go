package main

import (
	"groupietracker/web"
	"net/http"
)

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	web.UnMarshalData()
	web.DataFormat()
	web.SearchSuggest()
	web.PrintToConsole("Starting server at :8080")
	web.PrintToConsole("To stop, press CTRL+C")
	http.HandleFunc("/", web.HomePageHandler)
	http.HandleFunc("/artist/", web.ArtistHandler)
	http.HandleFunc("/search", web.SearchResultHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("error ListenAndServe related")
	}
}
