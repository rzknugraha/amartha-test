package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

// func init() {
// 	viper.AddConfigPath(".")
// 	viper.SetConfigName(".config")
// 	viper.SetConfigType("yaml")

// 	viper.AutomaticEnv() // read in environment variables that match
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file: " + viper.ConfigFileUsed())
// 	} else {
// 		log.Fatal(err)
// 	}
// 	err := viper.Unmarshal(&Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// type Config struct {
// 	Database struct {
// 		Driver   string `json:"driver"`
// 		Dsn      string `json:"dsn"`
// 		MaxConns int    `json:"max_cons"`
// 		MaxIdle  int    `json:"max_idle"`
// 	} `json:"database"`
// }

var validate *validator.Validate

//ShortenParam param for shor
type ShortenParam struct {
	URL       string `json:"url" validate:"required"`
	Shortcode string `json:"shortcode" validate:"omitempty,alphanum,len=6"`
}

//ShortenResp resp shorten
type ShortenResp struct {
	Shortcode string `json:"shortcode"`
	Message   string `json:"message"`
}

//DataURL struct for template data
type DataURL struct {
	Shortcode    string
	URL          string
	CreatedAt    string
	LastSeenDate string
	Count        int
}

//StatsResp for resp stats
type StatsResp struct {
	StartDate     string `json:"startDate" `
	LastSeenDate  string `json:"lastSeenDate"`
	RedirectCount int    `json:"redirectCount"`
}

//DataByShorten decalere mapping for temp data
var DataByShorten map[string][]*DataURL

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/shorten", short).Methods("POST")
	r.HandleFunc("/shortcode/{shorten}", getShort).Methods("GET")
	r.HandleFunc("/shortcode/stats/{shorten}", stats).Methods("GET")

	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	DataByShorten = make(map[string][]*DataURL)

	port := ":8008"
	fmt.Println("run in port " + port)
	log.Fatal(http.ListenAndServe(port, r))

}

//Response http response helper
func Response(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)

}

func short(w http.ResponseWriter, r *http.Request) {

	var p ShortenParam
	var dTemp DataURL
	var shorten string

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {

		Response(w, http.StatusInternalServerError, &ShortenResp{
			Message: "Error Decode",
		})
		return
	}

	err = validate.Struct(p)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {

			if e.StructField() == "Shortcode" {
				Response(w, http.StatusUnprocessableEntity, &ShortenResp{
					Message: "Error validation shortcode",
				})
				return
			}
		}
		Response(w, http.StatusBadRequest, &ShortenResp{
			Message: "Error validation",
		})
		return
	}

	if p.Shortcode == "" || len(p.Shortcode) < 0 {
		shorten = randSeq(6)
	} else {
		shorten = p.Shortcode
	}

	getData, _ := getDataByShorten(shorten)

	if getData {
		Response(w, http.StatusConflict, &ShortenResp{
			Message: "Exist Shorten",
		})

		return

	}

	dTemp.Shortcode = shorten
	dTemp.URL = p.URL
	dTemp.CreatedAt = time.Now().Format("2006-01-02T15:04:05+00:00")

	store(dTemp)

	Response(w, http.StatusCreated, &ShortenResp{
		Shortcode: shorten,
		Message:   "",
	})
	return

}

func getShort(w http.ResponseWriter, r *http.Request) {
	shorten := mux.Vars(r)["shorten"]

	if shorten == "" {
		Response(w, http.StatusNotFound, &ShortenResp{
			Message: "Not Found",
		})
		return
	}

	checkShorten, data := getDataByShorten(shorten)

	if checkShorten {

		update(data)

		http.Redirect(w, r, data.URL, http.StatusFound)
		return
	}
	Response(w, http.StatusNotFound, &ShortenResp{})
	return

}
func stats(w http.ResponseWriter, r *http.Request) {
	shorten := mux.Vars(r)["shorten"]

	if shorten == "" {
		Response(w, http.StatusNotFound, &ShortenResp{
			Message: "Not Found",
		})
		return
	}

	checkShorten, data := getDataByShorten(shorten)

	if checkShorten {
		resp := StatsResp{
			StartDate:     data.CreatedAt,
			LastSeenDate:  data.LastSeenDate,
			RedirectCount: data.Count,
		}
		Response(w, http.StatusOK, resp)
		return
	}
	Response(w, http.StatusNotFound, &ShortenResp{})
	return

}

func store(v DataURL) {

	DataByShorten[v.Shortcode] = append(DataByShorten[v.Shortcode], &v)
}
func update(v DataURL) {

	v.Count = v.Count + 1
	v.LastSeenDate = time.Now().Format(time.RFC3339)

	if len(DataByShorten[v.Shortcode]) > 0 {

		for _, val := range DataByShorten[v.Shortcode] {
			val.Count = v.Count
			val.LastSeenDate = v.LastSeenDate
		}
	}
}

func getDataByShorten(s string) (state bool, data DataURL) {

	if len(DataByShorten[s]) > 0 {

		for _, val := range DataByShorten[s] {

			data = DataURL{
				Shortcode:    val.Shortcode,
				URL:          val.URL,
				CreatedAt:    val.CreatedAt,
				LastSeenDate: val.LastSeenDate,
				Count:        val.Count,
			}
			return true, data

		}

	}

	return false, data
}
