package main

import (
	"be-tickitz/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type MovieListResponse struct {
	Results []struct {
		ID int `json:"id"`
	} `json:"results"`
}

type MovieDetail struct {
	Title        string `json:"title"`
	Overview     string `json:"overview"`
	ReleaseDate  string `json:"release_date"`
	Runtime      int    `json:"runtime"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
	Genres       []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

type Credits struct {
	Crew []struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
	Cast []struct {
		Name string `json:"name"`
	} `json:"cast"`
}

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func main() {
	godotenv.Load()
	apiKey := os.Getenv("TMDB_API_KEY")
	db, err := utils.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Release()

	categories := []string{"now_playing", "upcoming"}
	for _, cat := range categories {
		log.Println("📥 Fetching category:", cat)
		var list MovieListResponse
		if err := fetchJSON("https://api.themoviedb.org/3/movie/"+cat+"?api_key="+apiKey, &list); err != nil {
			log.Println("Failed to fetch", cat)
			continue
		}

		for _, item := range list.Results {
			var detail MovieDetail
			var credits Credits

			detailURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", item.ID, apiKey)
			creditURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s", item.ID, apiKey)

			if fetchJSON(detailURL, &detail) != nil || fetchJSON(creditURL, &credits) != nil {
				continue
			}

			var movieID int
			err := db.QueryRow(context.Background(), `
				INSERT INTO movies (title, description, release_date, duration_minutes, image, horizontal_image)
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
			`, detail.Title, detail.Overview, detail.ReleaseDate, detail.Runtime,
				"https://image.tmdb.org/t/p/w500"+detail.PosterPath,
				"https://image.tmdb.org/t/p/original"+detail.BackdropPath,
			).Scan(&movieID)
			if err != nil {
				log.Println("❌ Insert movie failed:", err)
				continue
			}

			for _, g := range detail.Genres {
				var genreID int
				db.QueryRow(context.Background(), `
					INSERT INTO genres (genre_name) VALUES ($1)
					ON CONFLICT (genre_name) DO UPDATE SET genre_name = EXCLUDED.genre_name
					RETURNING id
				`, g.Name).Scan(&genreID)

				db.Exec(context.Background(), `INSERT INTO movie_genres (id_movie, id_genre) VALUES ($1, $2)`, movieID, genreID)
			}

			for _, crew := range credits.Crew {
				if crew.Job == "Director" {
					var directorID int
					db.QueryRow(context.Background(), `
						INSERT INTO directors (director_name) VALUES ($1)
						ON CONFLICT (director_name) DO UPDATE SET director_name = EXCLUDED.director_name
						RETURNING id
					`, crew.Name).Scan(&directorID)

					db.Exec(context.Background(), `INSERT INTO movie_directors (id_movie, id_director) VALUES ($1, $2)`, movieID, directorID)
					break
				}
			}

			for i, cast := range credits.Cast {
				if i >= 6 {
					break
				}
				var actorID int
				db.QueryRow(context.Background(), `
					INSERT INTO actors (actor_name) VALUES ($1)
					ON CONFLICT (actor_name) DO UPDATE SET actor_name = EXCLUDED.actor_name
					RETURNING id
				`, cast.Name).Scan(&actorID)

				db.Exec(context.Background(), `INSERT INTO movie_casts (id_movie, id_actor) VALUES ($1, $2)`, movieID, actorID)
			}

			log.Println("✅ Inserted:", detail.Title)
		}
	}
	log.Println("🎉 Done seeding all movies with genres, directors, and casts")
}
