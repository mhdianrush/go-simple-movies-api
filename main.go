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
	// using the json encoder ==> to show to user that the all of the moviess
}

// Delete Function
func deleteMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	// using params ==> to detect the specific id (by json format) by getting the item.ID below
	for i, item := range movies {
		if item.Id == params["id"] {
			// delete the movie by index
			movies = append(movies[:i], movies[i+1:]...)
			// the movie was deleted, and the now we only have the movie that not deleted

			// don't use break to stop the delete operation of the index, or then will deleted only one by one
		}
	}
	// after delete the movie, show again the newest movies that exist to all of the user
	json.NewEncoder(writer).Encode(movies)
}

// Get By Id Function
func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	// using params ==> to detect the specific id (by json format) by getting the item.ID below
	for _, item := range movies {
		// we don't need the index beacuse we have to only show the movies by spesific id
		if item.Id == params["id"] {
			json.NewEncoder(writer).Encode(item)
			// item inside the Encode is used to select only one the movie that relate to the id

			// don't use return in the end, then, Get By Id appear only one at a time
		}
	}
}

// Create Function
func createMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var movie Movie
	// not use the pointer to data type Movie because the movie variable only follow the struct Movie contract

	// because of create movie, we will send something to the body (request to body)
	json.NewDecoder(request.Body).Decode(&movie)
	// the movie was enter to the temporary variable

	// now determine the id of the new movie, beacuse we don't construct the database that makes auto_increment
	movie.Id = strconv.Itoa(rand.Intn(1000000000))
	// use package math/rand and inside the package math/rand use the Intn function to catch the int value
	// 10 inside the Intn above is only example value, we can use others number. up to you.
	// we have to convert this to string format because in the Movie sruct above, we declare that the id is string
	movies = append(movies, movie)
	json.NewEncoder(writer).Encode(movie)
	// to show the movie that created
}

// Update Function
func updateMovie(writer http.ResponseWriter, request *http.Request) {
	// especially for update movie, first, we have to delete the specific movie by id, then we make a new movie by the specific id too
	/**
	set json content type
	params
	loop over the movie ==> using for range
	delete the movie with the specified id that we've sent
	add a new movie ==> the movie that we send in the body
	*/

	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	// for the first, delete by id
	for i, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
		}
	}
	// then, add new data
	var movie Movie
	json.NewDecoder(request.Body).Decode(&movie)
	movie.Id = params["id"]
	// "id" inside params equal to id in the URL parameter
	// in this new movie id, we still use the the same id with the earlier movie that deleted
	movies = append(movies, movie)
	// show movie to the client
	json.NewEncoder(writer).Encode(movie)
	// using movie inside the Encode because we want to show the newest movie that added
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

	// or we can use ==> log.Fatal(http.ListenAndServe(":8000", r))
}
