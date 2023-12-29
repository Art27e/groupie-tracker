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

	tmpl, _ := template.New("index").Funcs(template.FuncMap{"jsEscape": JsEscape}).Parse(`<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" type="text/css" href="/css/style.css">
        <link href="https://fonts.cdnfonts.com/css/circular-std"
            rel="stylesheet">
        <link href='https://unpkg.com/boxicons@2.1.4/css/boxicons.min.css'
            rel='stylesheet'>
        <link rel="icon"
            href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAB4AAAAeCAYAAAA7MK6iAAAAAXNSR0IArs4c6QAAAYJJREFUSEvllf0xRDEUR89WgA50gArQARWgElSCCugAFaCCLQEdmDOTmGy8vNw3+8aOcf/afcnNuV/5ZcGGbLEhLn8afA4sgecp1Zsj4w9gC7gDDCJk64L3gZdEugBuQ1RYu8dmeJNgB8Drb4HN8Az4BLajUPetW2oz3EuDdTQ32HIeAvawNDN8Tx+ugav02yo89frdyziX0jPr4THDxwQ7TrCy5/rWwX4H3gKbjYc6tZo9FFQOjxlepvUdwGvlfrP1emnuNyjXVmwIrLPQPCxvCVo7C7AFrucAPVw/1+y9pp/wlYmvwWWpdBoThZ5wlG360aoSfALcF/UYE4SocNSJnAIP9XWaAo4KRwhsIPXG1mRGhCNc6lzlyGSOCUfkRjSVa2wyW8Jh4NEb0ZXMIQEZEo5crbJVo89kT7ly34Xlt3ZIOEpxmEUya8Hxf0s4hvY2v0Uyrp17whEKYCo4Khxd+FRwVDhmB0eEowutJTPkkF6q3d5D3ztsaql754XX/x/4C7SvZR+uZb5OAAAAAElFTkSuQmCC"
            type="image/x-icon">
        <script src="https://api-maps.yandex.ru/2.1/?apikey=ab421167-6629-47f8-b0b0-383cb58a54c7&lang=en_US" type="text/javascript"></script>
        <script type="text/javascript">
		var concertsMapData = {{.ConcertsMap}};

            ymaps.ready(init);
            function init(){
                var map = new ymaps.Map("map", {
                    center: [42.09, -55.019], // focus on atlantic ocean
                    zoom: 2
                });
				// Add placemarks to the map
				for (var city in concertsMapData) {
					if (concertsMapData.hasOwnProperty(city)) {
						var coordinatesString = concertsMapData[city];
						var coordinatesArray = coordinatesString.split(/[:\s]+/).map(function(coord) {
							return parseFloat(coord);
						});
		
						var placemark = new ymaps.Placemark([coordinatesArray[1], coordinatesArray[0]], { hintContent: city }); // обратите внимание на порядок координат
						map.geoObjects.add(placemark);
					}
				}
			}
        </script>
        <title>Groupie Tracker</title>
    </head>
    <body>
        <div class="logo-container">
            <div class="logo">
                <i class='bx bxs-music'></i>
            </div>
            <div class="header">
                <h1 class="page-header">Groupie-Tracker</h1>
            </div>
        </div>
        <br>
        <button type="button" onclick="history.back()">GO BACK</button>
        <div class="row">
            <div class="group-card-info">
                <h1 class=infopage> {{ .Name }} </h1>
                <img class="infocard-image" src="{{.Image}}"></img>
            <header>Concert dates</header>
            {{range $key, $value := .Relation}}
            <p class="info-card-city">{{ $key }}</p>
            {{ range . }}
            <p class="info-card-date">{{ . }}</p>
            {{ end }}
            {{end}}
        </div>
		<div id="map" style="width:600px; height: 400px; margin-left: 10px; margin-top: 50px; ">
        </div>
    </div>
</body>`)

	idInt := 1

	id := strings.TrimPrefix(request.URL.Path, "/artist/")
	idInt, err := strconv.Atoi(id)

	MakeCoords(idInt - 1)

	if len(request.URL.Path) > 8 {
		if err != nil {
			http.Error(write, "Oops, something went wrong", 500)
			return
		}
		if idInt > 52 {
			http.Error(write, "Oops, something went wrong", 500)
			return
		} else {
			err = tmpl.Execute(write, data.HomePageData.ArtistInfo[idInt-1])
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
