package web

import (
	"encoding/json"
	"fmt"
	"groupietracker/data"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func getAPI(source string) string {
	response, err := http.Get(source)
	if err != nil {
		panic("api is not found")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("problem with api reading")
	}
	return string(body)
}

func UnMarshalData() {
	if json.Valid([]byte(getAPI("https://groupietrackers.herokuapp.com/api"))) {
		if err := json.Unmarshal([]byte(getAPI("https://groupietrackers.herokuapp.com/api")), &data.GetWebApiData); err != nil {
			panic(err)
		}
	}
	if json.Valid([]byte(getAPI(data.GetWebApiData.Artists))) {
		if err := json.Unmarshal([]byte(getAPI(data.GetWebApiData.Artists)), &data.ArtistInfo); err != nil {
			panic(err)
		}
	}
	if json.Valid([]byte(getAPI(data.GetWebApiData.Relation))) {
		if err := json.Unmarshal([]byte(getAPI(data.GetWebApiData.Relation)), &data.RelationInfo); err != nil {
			panic(err)
		}
	}
	PrintToConsole("All APIs are succesfully loaded")
}

func MatchResult(collection interface{}, search string) bool {
	switch reflect.TypeOf(collection).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(collection)
		for i := 0; i < s.Len(); i++ {
			if strings.EqualFold(s.Index(i).String(), search) {
				return true
			}
		}
	case reflect.Map:
		m := reflect.ValueOf(collection)
		for _, v := range m.MapKeys() {
			if strings.EqualFold(v.String(), search) {
				return true
			}
		}
	}
	return false
}

func DataFormat() {
	// Patterns for formatting
	patterntoUcFirstLet := regexp.MustCompile(`^[a-z]`)
	patternForDeph := regexp.MustCompile(`-`)
	patternForSpc := regexp.MustCompile(`_`)
	patternToUc := regexp.MustCompile(`(?:^|\,\s)([a-z])`)
	patternToUsUk := regexp.MustCompile(`Usa|Uk`)
	patternToTitle := regexp.MustCompile(`\s([a-z])`)

	for x := range data.ArtistInfo {
		storeData := []string{}
		tempMap := make(map[string][]string)
		for key := range data.RelationInfo.Index[x].Relations {
			newKey := patterntoUcFirstLet.ReplaceAllStringFunc(key, func(s string) string {
				return strings.ToUpper(s)
			})
			newKey = patternForDeph.ReplaceAllString(newKey, ", ")
			newKey = patternForSpc.ReplaceAllString(newKey, " ")
			newKey = patternToUc.ReplaceAllStringFunc(newKey, func(s string) string {
				return strings.ToUpper(string(s))
			})
			newKey = patternToUsUk.ReplaceAllStringFunc(newKey, func(s string) string {
				return strings.ToUpper(string(s[0:]))
			})
			newKey = patternToTitle.ReplaceAllStringFunc(newKey, func(s string) string {
				return strings.ToUpper(string(s))
			})
			storeData = nil
			for i := range data.RelationInfo.Index[x].Relations[key] {
				storeData = append(storeData, data.RelationInfo.Index[x].Relations[key][i])
			}
			tempMap[newKey] = storeData
		}
		data.ArtistInfo[x].Relation = tempMap
	}
}

func SearchSuggest() {
	resDates := []int{}
	resLocs := []string{}
	usedDates := make(map[int]bool)
	usedLocations := make(map[string]bool)

	for i, a := range data.ArtistInfo { // res slice is without duplicates
		if !usedDates[a.Date] {
			resDates = append(resDates, a.Date)
			usedDates[a.Date] = true
		}
		for key := range data.ArtistInfo[i].Relation {
			if !usedLocations[key] {
				resLocs = append(resLocs, key)
				usedLocations[key] = true
			}
		}
	}

	data.Searching = data.SearchInfo{
		Date: resDates,
		Locs: resLocs,
	}
	data.HomePageData = data.MainData{
		ArtistInfo: data.ArtistInfo,
		Search:     data.Searching,
	}
}

func PrintToConsole(s string) {
	toBytes := []byte(s)
	for i := 0; i < len(toBytes); i++ {
		if toBytes[i] >= 32 && toBytes[i] <= 126 {
			fmt.Print(string(s[i]))
			time.Sleep(time.Second / 15)
		} else {
			panic("unexpected error")
		}
	}
	fmt.Println()
}
