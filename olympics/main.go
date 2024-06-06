//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
)

/*--------------structs-----------------*/
type Athlete struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type Medals struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type AthleteInfo struct {
	Athlete      string          `json:"athlete"`
	Country      string          `json:"country"`
	Medals       Medals          `json:"medals"`
	MedalsByYear map[int]*Medals `json:"medals_by_year"`
}

type CountryInfo struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

/*--------------information_about_athlete-----------------*/

func atletInformation(a *Athlete) AthleteInfo {
	return AthleteInfo{Athlete: a.Athlete,
		MedalsByYear: make(map[int]*Medals),
		Country:      a.Country,
	}
}

/*--------------filter-----------------*/

func Filter(atletes []Athlete, func_filter func(Athlete) bool) (result []Athlete) {
	for _, atlete := range atletes {
		if func_filter(atlete) {
			result = append(result, atlete)
		}
	}
	return
}

func FillMedals(info *AthleteInfo, atlete Athlete) {
	medals, ok := info.MedalsByYear[atlete.Year]
	if !ok {
		medals = &Medals{0, 0, 0, 0}
		info.MedalsByYear[atlete.Year] = medals
	}
	totalMedals := atlete.Gold + atlete.Bronze + atlete.Silver

	info.Medals.Bronze += atlete.Bronze
	info.Medals.Silver += atlete.Silver
	info.Medals.Gold += atlete.Gold
	info.Medals.Total += totalMedals

	info.MedalsByYear[atlete.Year].Bronze += atlete.Bronze
	info.MedalsByYear[atlete.Year].Silver += atlete.Silver
	info.MedalsByYear[atlete.Year].Gold += atlete.Gold
	info.MedalsByYear[atlete.Year].Total += totalMedals
}

func GetAllAthlets(atlets []Athlete) map[string]*AthleteInfo {
	allAtlets := make(map[string]*AthleteInfo)
	for _, atlete := range atlets {
		info, ok := allAtlets[atlete.Athlete]
		if !ok {
			info = &AthleteInfo{
				Athlete:      atlete.Athlete,
				MedalsByYear: make(map[int]*Medals),
				Country:      atlete.Country,
			}
			allAtlets[atlete.Athlete] = info
		}

		FillMedals(info, atlete)
	}
	return allAtlets
}

func GetCountries(atlets []Athlete) map[string]*CountryInfo {
	allAtlets := make(map[string]*CountryInfo)
	for _, atlete := range atlets {
		if _, ok := allAtlets[atlete.Country]; !ok {
			allAtlets[atlete.Country] = &CountryInfo{
				Country: atlete.Country,
				Gold:    0,
				Silver:  0,
				Bronze:  0,
			}
		}
		totalMedals := atlete.Gold + atlete.Bronze + atlete.Silver
		allAtlets[atlete.Country].Bronze += atlete.Bronze
		allAtlets[atlete.Country].Silver += atlete.Silver
		allAtlets[atlete.Country].Gold += atlete.Gold
		allAtlets[atlete.Country].Total += totalMedals
	}
	return allAtlets
}

var athletes []Athlete

func WriteJson(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		http.Error(w, "error write", http.StatusBadRequest)
	}
}

func GetSortedInfo(filteredAthlets map[string]*AthleteInfo) []*AthleteInfo {
	values := make([]*AthleteInfo, 0, len(filteredAthlets))
	for _, v := range filteredAthlets {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Medals.Gold != values[j].Medals.Gold {
			return values[i].Medals.Gold > values[j].Medals.Gold
		}
		if values[i].Medals.Silver != values[j].Medals.Silver {
			return values[i].Medals.Silver > values[j].Medals.Silver
		}
		if values[i].Medals.Bronze != values[j].Medals.Bronze {
			return values[i].Medals.Bronze > values[j].Medals.Bronze
		}

		return values[i].Athlete < values[j].Athlete
	})
	return values
}

func GetSortedCountry(filteredAthlets map[string]*CountryInfo) []*CountryInfo {
	values := make([]*CountryInfo, 0, len(filteredAthlets))
	for _, v := range filteredAthlets {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Gold != values[j].Gold {
			return values[i].Gold > values[j].Gold
		}
		if values[i].Silver != values[j].Silver {
			return values[i].Silver > values[j].Silver
		}
		if values[i].Bronze != values[j].Bronze {
			return values[i].Bronze > values[j].Bronze
		}

		return values[i].Country < values[j].Country
	})
	return values
}

func AthleteInfoHundler(w http.ResponseWriter, r *http.Request) {
	var name string
	query := r.URL.Query()
	// проверим есть ли вообще имя
	nameQuery, ok := query["name"]
	if !ok || len(nameQuery[0]) == 0 {
		http.Error(w, "no name", http.StatusBadRequest)
		return
	}
	name = nameQuery[0]

	filter := Filter(athletes, func(a Athlete) bool {
		return a.Athlete == name
	})
	if len(filter) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	responce := atletInformation(&filter[0])
	for _, athlete := range filter {
		FillMedals(&responce, athlete)
	}
	jsonMarshal, err := json.Marshal(&responce)
	if err != nil {
		http.Error(w, "error marshal json", http.StatusBadRequest)
		return
	}
	WriteJson(w, r, jsonMarshal)

}

func TopAthletesInSportHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	queries := r.URL.Query()
	sport := queries.Get("sport")
	limitStr := queries.Get("limit")

	// Проверяем наличие параметра "sport"
	if sport == "" {
		http.Error(w, "no sport param", http.StatusBadRequest)
		return
	}

	// Парсим параметр "limit"
	limit := 3 // значение по умолчанию
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit param", http.StatusBadRequest)
			return
		}
	}

	// Фильтруем атлетов по указанному виду спорта
	filteredAthletes := Filter(athletes, func(a Athlete) bool {
		return a.Sport == sport
	})

	// Если не найдены атлеты по указанному виду спорта, возвращаем ошибку
	if len(filteredAthletes) == 0 {
		http.Error(w, "not found sport", http.StatusNotFound)
		return
	}

	// Получаем информацию о фильтрованных атлетах
	filteredAthletesInfo := GetAllAthlets(filteredAthletes)

	// Сортируем информацию о атлетах по количеству медалей
	sortedAthletesInfo := GetSortedInfo(filteredAthletesInfo)

	// Ограничиваем количество возвращаемых результатов
	limits := int(math.Min(float64(limit), float64(len(sortedAthletesInfo))))
	result := sortedAthletesInfo[:limits]

	// Преобразуем результат в формат JSON
	jsonResult, err := json.Marshal(&result)
	if err != nil {
		http.Error(w, "marshal error", http.StatusBadRequest)
		return
	}

	// Отправляем JSON-ответ клиенту
	WriteJson(w, r, jsonResult)
}

func TopCountriesInYearHandler(w http.ResponseWriter, r *http.Request) {
	// Инициализация переменных
	var yearParam int
	var limit int
	var err error

	// Получение параметров запроса
	queries := r.URL.Query()

	// Парсинг параметра "year"
	yearSlice, ok := queries["year"]
	if !ok {
		http.Error(w, "no year param", http.StatusBadRequest)
		return
	}
	if yearSlice[0] == "" {
		http.Error(w, "empty year param", http.StatusBadRequest)
		return
	}
	yearParam, err = strconv.Atoi(yearSlice[0])
	if err != nil {
		http.Error(w, "invalid year param", http.StatusBadRequest)
		return
	}

	// Парсинг параметра "limit"
	limitSlice, ok := queries["limit"]
	if !ok || limitSlice[0] == "" {
		limit = 3 // Значение по умолчанию
	} else {
		limit, err = strconv.Atoi(limitSlice[0])
		if err != nil {
			http.Error(w, "invalid limit param", http.StatusBadRequest)
			return
		}
	}

	// Фильтрация атлетов по указанному году
	filtered := Filter(athletes, func(a Athlete) bool {
		return a.Year == yearParam
	})

	// Если атлетов не найдено, возвращаем ошибку
	if len(filtered) == 0 {
		http.Error(w, "no athletes found for the year", http.StatusNotFound)
		return
	}

	// Получение информации о странах с учетом отфильтрованных атлетов
	filteredCountries := GetCountries(filtered)
	sortedCountries := GetSortedCountry(filteredCountries)

	// Ограничение количества возвращаемых результатов
	limits := int(math.Min(float64(limit), float64(len(sortedCountries))))
	result := sortedCountries[:limits]

	// Преобразование результата в JSON и отправка клиенту
	jsonResult, err := json.Marshal(&result)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	WriteJson(w, r, jsonResult)
}
func main() {
	port := flag.String("port", "80", "http server port")

	path := flag.String("data", "./olympics/testdata/olympicWinners.json", "path json")
	flag.Parse()
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal("reading error")
	}
	value, _ := ioutil.ReadAll(file)
	file.Close()
	err = json.Unmarshal(value, &athletes)
	if err != nil {
		log.Fatal("unmarshall fatal")
	}
	http.HandleFunc("/athlete-info", AthleteInfoHundler)
	http.HandleFunc("/top-athletes-in-sport", TopAthletesInSportHandler)
	http.HandleFunc("/top-countries-in-year", TopCountriesInYearHandler)
	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
