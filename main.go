package main

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nalgeon/redka"
	"github.com/torbatti/sqlite-cache/models"
	"github.com/torbatti/sqlite-cache/utils"
	"github.com/torbatti/sqlite-cache/views"

	_ "embed"
)

var json_paths = []string{
	"json/3DS.json",
	"json/Dreamcast.json",
	"json/DS.json",
	"json/Gamecube.json",
	"json/GB.json",
	"json/GBA.json",
	"json/GBC.json",
	"json/N64.json",
	"json/NES.json",
	"json/PS1.json",
	"json/PS2.json",
	"json/PS3.json",
	"json/SegaGenesis.json",
	"json/SegaSaturn.json",
	"json/SNES.json",
	"json/Wii.json",
	"json/XBOX.json",
	"json/XBOX360.json",
}

//go:embed models/sqlc/schema/games.sql
var ddl string

type App struct {
	DB     *sql.DB
	RK     *redka.DB
	DB_CTX context.Context
}

func init_db() (*sql.DB, context.Context) {

	ctx := context.Background()

	// db, err := sql.Open("sqlite3", ":memory:")
	// db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory&_journal=WAL")
	db, err := sql.Open("sqlite3", "file:test.db?_journal=WAL")
	if err != nil {
		log.Fatal(err)
	}

	queries := models.New(db)

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		println("Creating Table Error :", err)
	}

	// check for an empty struct
	// if game != (models.Game{}) {
	// }

	// // list all authors
	// games, err := queries.ListGames(ctx)
	// if err != nil {
	// 	println(err)
	// }
	// log.Println(games)
	_, err = queries.GetGame(ctx, 1)
	if err == sql.ErrNoRows {
		for _, path := range json_paths {
			var games = utils.OpenParseJson(path)

			for _, game := range games {
				_, err := queries.CreateGame(ctx, models.CreateGameParams{
					Game:      game.Game,
					Year:      game.Year,
					Dev:       sql.NullString{String: game.Dev}.String,
					Publisher: sql.NullString{String: game.Publisher}.String,
					Platform:  sql.NullString{String: game.Platform}.String,
				})
				if err != nil {
					println(err)
				}
			}
		}
	}
	if err != sql.ErrNoRows && err != nil {
		panic(err)
	}

	return db, ctx
}

func init_rk() *redka.DB {
	rk, err := redka.Open("data.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Always close the database when you are finished.
	return rk
}

func main() {

	db, ctx := init_db()
	rk := init_rk()
	defer rk.Close()

	r := init_routes(db, ctx, rk)

	server := &http.Server{
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		// WriteTimeout: 60 * time.Second, // breaks sse!
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}
	println("Starting server at : ", server.Addr)
	server.ListenAndServe()
}

func init_routes(db *sql.DB, ctx context.Context, rk *redka.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "text/html", "text/css"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	r.Get("/game/{gameId}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "gameId")
		// println(param)

		cached_html, err := rk.Str().Get(param)
		if err != nil {
			println(err)
		}
		// println(cached_html.String())

		if cached_html.String() == "" {
			println("from db")

			i, err := strconv.Atoi(param)
			if err != nil {
				println(err)
			}

			queries := models.New(db)
			game, _ := queries.GetGame(ctx, int64(i))

			info := views.PageInfo{
				HeadInfo: views.HeadInfo{
					Title:       "تست",
					Description: "تست سیکولایت و کش",
				},

				Game: game,
			}

			var buffer = bytes.NewBufferString("")

			views.BaziPage(info).Render(r.Context(), buffer)
			rk.Str().Set(param, buffer.String())

			w.Write([]byte(buffer.String()))
		} else {
			w.Write([]byte(cached_html.String()))
			println("from cache")
		}

	})
	return r
}
