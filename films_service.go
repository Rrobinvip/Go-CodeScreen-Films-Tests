package films

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"
)

const filmsEndpointUrl string = "https://toolbox.palette-adv.spectrocloud.com:5002/films"

// Your API token. Needed to successfully authenticate when calling the films endpoint.
// Must be included in the Authorization header in the request sent to the films endpoint.
const apiToken string = "8c5996d5-fb89-46c9-8821-7063cfbc18b1"

type Film struct {
	Name         string  `json:"name"`
	Length       int     `json:"length"`
	Rating       float64 `json:"rating"`
	ReleaseDate  string  `json:"releaseDate"`
	DirectorName string  `json:"directorName"`
}

// Added singleton design, to make api call to film for only once
type FilmSingleton struct {
	films []Film
	once  sync.Once
}

var filmInstance *FilmSingleton
var once sync.Once

func GetFilmInstance() *FilmSingleton {
	once.Do(func() {
		filmInstance = &FilmSingleton{}
		filmInstance.LoadFilms()
	})

	return filmInstance
}

func (fi *FilmSingleton) LoadFilms() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", filmsEndpointUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Something wrong with the server")
		return
	}

	respBodyRow, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var films []Film
	json.Unmarshal(respBodyRow, &films)

	fi.films = films
}

// GetFilms retrieves the data for all films by calling the https://toolbox.palette-adv.spectrocloud.com:5002/films endpoint.
func GetFilms() []Film {
	//TODO Implement
	return GetFilmInstance().films
}

// BestRatedFilm retrieves the name of the best rated film that was directed by the director with the given name.
// If there are no films directed by the given director, return an empty string.
// Note: there will only be one film with the best rating.
func BestRatedFilm(directorName string) string {
	//TODO Implement
	films := GetFilms()

	res := ""
	currentHighest := 0

	for _, f := range films {
		if f.DirectorName == directorName {
			if f.Rating > float64(currentHighest) {
				res = f.Name
				currentHighest = int(f.Rating)
			}
		}
	}

	return res
}

// DirectorWithMostFilms retrieves the name of the director who has directed the most films
// in the CodeScreen Film service.
func DirectorWithMostFilms() string {
	//TODO Implement
	films := GetFilms()
	res := ""

	// As I mentioned during the interview, I would like to create a hashmap to store which director has how many films
	directorCount := make(map[string]int)
	for _, f := range films {
		directorCount[f.DirectorName]++
	}

	currentHighest := 0

	// Loop through the hashmap and find out the highest!
	for director, count := range directorCount {
		if count > currentHighest {
			res = director
			// I remeber to update the currentHighest this time :D
			currentHighest = count
		}
	}

	return res
}

// AverageRating retrieves the average rating for the films directed by the given director, rounded to 1 decimal place.
// If there are no films directed by the given director, return 0.0.
func AverageRating(directorName string) float64 {
	//TODO Implement
	films := GetFilms()
	ratingSum := 0.0
	count := 0

	for _, f := range films {
		if f.DirectorName == directorName {
			ratingSum += f.Rating
			count += 1
		}
	}

	// I like to do edge case checks to handle error, if the director is not given, I will simply return 0
	if ratingSum == 0.0 {
		return 0.0
	} else {
		avg := ratingSum / float64(count)
		return math.Round(avg*10) / 10
	}
}

/*
ShortestFilmReleaseGap retrieves the shortest number of days between any two film releases directed by the given director.
If there are no films directed by the given director, return 0.
If there is only one film directed by the given director, return 0.
Note: no director released more than one film on any given day.

For example, if the service returns the following 3 films:

	{
	    "name": "Batman Begins",
	    "length": 140,
	    "rating": 8.2,
	    "releaseDate": "2006-06-16",
	    "directorName": "Christopher Nolan"
	},

	{
	    "name": "Interstellar",
	    "length": 169,
	    "rating": 8.6,
	    "releaseDate": "2014-11-07",
	    "directorName": "Christopher Nolan"
	},

	{
	    "name": "Prestige",
	    "length": 130,
	    "rating": 8.5,
	    "releaseDate": "2006-11-10",
	    "directorName": "Christopher Nolan"
	}

Then this method should return 147 for Christopher Nolan, as Prestige was released 147 days after Batman Begins.
*/
func ShortestFilmReleaseGap(directorName string) int {
	//TODO Implement
	films := GetFilms()
	var releaseDates []time.Time

	// This is where I need to check the time. As I mentioned at the beginging of the interview, I would like to set the ReleaseDate to a time.Time
	// so it can be recognized when data is binding. But I can convert it by now.
	for _, f := range films {
		if f.DirectorName == directorName {
			releaseDate, err := time.Parse("2006-01-02", f.ReleaseDate)
			if err != nil {
				fmt.Println(err)
				continue
			}
			releaseDates = append(releaseDates, releaseDate)
		}
	}

	if len(releaseDates) < 2 {
		return 0
	}

	// I use a lambda function to sort it.
	sort.Slice(releaseDates, func(i, j int) bool {
		return releaseDates[i].Before(releaseDates[j])
	})

	res := math.MaxInt64
	for i := 0; i < len(releaseDates)-1; i++ {
		// I did a little bit of search here to figure out how to calculate days between two dates.
		gap := int((releaseDates[i+1].Sub(releaseDates[i]).Hours()) / 24)
		if gap < res {
			res = gap
		}
	}

	return res
}

// Now it can pass all the test!
