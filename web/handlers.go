package web

import (
	"groupietracker/data"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func HomePageHandler(write http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		html, err := template.ParseFiles("template/index.html")
		write.Header().Set("Content-Type", "text/html")
		write.WriteHeader(http.StatusOK)
		if err != nil {
			http.Error(write, "Oops, something went wrong", 500)
			return
		}
		err = html.Execute(write, data.HomePageData)
		if err != nil {
			http.Error(write, "Oops, something went wrong", 500)
			return
		}
	}

	switch request.Method {
	case "GET":
		if request.URL.Path != "/" {
			http.Error(write, "Oops, something went wrong", 404)
			return
		}
	}
}

func ArtistHandler(write http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/info.html")
	if err != nil {
		http.Error(write, "Oops, something went wrong", 500)
		return
	}

	idInt := 1
	if len(request.URL.Path) > 8 {
		id := strings.TrimPrefix(request.URL.Path, "/artist/")
		idInt, err = strconv.Atoi(id)
		if err != nil {
			http.Error(write, "Oops, something went wrong", 500)
			return
		}
		if idInt > 52 {
			http.Error(write, "Oops, something went wrong", 500)
			return
		} else {
			err = html.Execute(write, data.HomePageData.ArtistInfo[idInt-1])
			if err != nil {
				http.Error(write, "Oops, something went wrong", 500)
				return
			}
		}
	} else {
		http.Error(write, "Oops, something went wrong", 404)
		return
	}
}

func SearchResultHandler(write http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/search.html")
	if err != nil {
		http.Error(write, "Oops, something went wrong PARSING", 500)
		return
	}
	if request.URL.Path != "/search" {
		http.Error(write, "Oops, something went wrong PATH", 500)
		return
	}
	searchResult := request.FormValue("artist")
	data.SearchData = nil // otherwise it will remember every search query

	// match search results with database
	for _, a := range data.ArtistInfo {
		if a.Name == searchResult || MatchResult(a.Members, searchResult) || strconv.Itoa(a.Date) == searchResult || a.FirstAlbum == searchResult || MatchResult(a.Relation, searchResult) {
			data.SearchData = append(data.SearchData, a)
			continue
		}
	}
	err = html.Execute(write, data.SearchData)
	if err != nil {
		http.Error(write, "Oops, something went wrong. EXECUTE", 500)
		return
	}
}
