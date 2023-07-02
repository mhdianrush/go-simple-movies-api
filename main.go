package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// every movie has one director

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

// Get All Function
func getMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(movies)
}

// Delete Function
func deleteMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	
	for i, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
		}
	}
	json.NewEncoder(writer).Encode(movies)
}

// Get By Id Function
func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(writer).Encode(item)
		}
	}
}

// Create Function
func createMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var movie Movie
	
	json.NewDecoder(request.Body).Decode(&movie)
	
	movie.Id = strconv.Itoa(rand.Intn(1000000000))
	
	movies = append(movies, movie)
	json.NewEncoder(writer).Encode(movie)
}

// Update Function
func updateMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	for i, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
		}
	}
	var movie Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movie.Id = params["id"]
	movies = append(movies, movie)
	json.NewEncoder(writer).Encode(movie)
}

func main() {
	r := mux.NewRouter()
	logger := logrus.New()

	movies = append(movies, Movie{
		Id:    "1",
		Isbn:  "12345",
		Title: "Movie One",
		Director: &Director{
			Firstname: "Muhammad",
			Lastname:  "Ian",
		},
	})
	movies = append(movies, Movie{
		Id:    "2",
		Isbn:  "12346",
		Title: "Movie Two",
		Director: &Director{
			Firstname: "Muhammad",
			Lastname:  "Rush",
		},
	})
	movies = append(movies, Movie{
		Id:    "3",
		Isbn:  "12347",
		Title: "Movie Three",
		Director: &Director{
			Firstname: "Ian",
			Lastname:  "Rush",
		},
	})

	// Get All
	r.HandleFunc("/movies", getMovies).Methods("GET")
	// Get By Id
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	// Create
	r.HandleFunc("/movies", createMovie).Methods("POST")
	// Update
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	// Delete
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting Server at Port 8000")

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: r,
	}
	err := server.ListenAndServe()
	logger.Fatal(err)
}
