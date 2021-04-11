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

type ShortenParam struct {
	URL       string `json:"url" validate:"required"`
	Shortcode string `json:"shortcode" validate:"omitempty,alphanum,len=6"`
}

type ShortenResp struct {
	Shortcode string `json:"shortcode"`
	Message   string `json:"message"`
}

type DataURL struct {
	Shortcode string
	URL       string
}

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
		fmt.Println(err)
		Response(w, http.StatusInternalServerError, &ShortenResp{
			Message: "Error Decode",
		})
		return
	}
	fmt.Println(p)

	err = validate.Struct(p)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			fmt.Println(e)
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

	if p.Shortcode == "" {
		shorten = randSeq(6)
	} else {
		shorten = p.Shortcode
	}
	fmt.Println("shorten")
	fmt.Println(shorten)
	getData := getDataByShorten(shorten)

	if getData {
		Response(w, http.StatusConflict, &ShortenResp{
			Message: "Exist Shorten",
		})

		return

	}
	dTemp.Shortcode = p.Shortcode
	dTemp.URL = p.URL

	store(dTemp)

	Response(w, http.StatusCreated, &ShortenResp{
		Shortcode: shorten,
		Message:   "",
	})
	return

}

func store(v DataURL) {

	DataByShorten[v.Shortcode] = append(DataByShorten[v.Shortcode], &v)
}

func getDataByShorten(s string) bool {
	fmt.Println("s")
	fmt.Println(s)
	if len(DataByShorten[s]) > 0 {
		for _, val := range DataByShorten[s] {
			fmt.Println("val")
			fmt.Println(val)
			if val.Shortcode != "" {
				return true
			}
		}
	}

	return false
}
