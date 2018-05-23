package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"metronom/password"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s := &server{}
	s.setup()

	s.ListenAndServe()
}

type server struct {
	http.Server
	router *chi.Mux
}

func (s *server) setup() {
	viper.SetDefault("HTTPaddr", ":8080")
	viper.AutomaticEnv()

	s.Addr = viper.GetString("HTTPaddr")

	s.router = chi.NewRouter()
	s.commonMiddleware()
	s.routes()
	s.Handler = s.router
}

func (s *server) routes() {
	s.router.Get("/v1/{amount}", s.handleREST())
	s.router.Get("/v1", s.handleREST())
}

func (s *server) commonMiddleware() {
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middleware.DefaultCompress)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
}

func (s *server) handleREST() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := extractParameter(r)

		err := password.Validate(p)
		if err != nil {
			log.Println(err)
			if err != password.ErrParametersExceedMinLength {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		p = password.AutoComplete(p)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		amount := chi.URLParam(r, "amount")
		a, err := strconv.Atoi(amount)
		if err != nil {
			a = 1
		}

		for i := 0; i < a; i++ {
			if r.FormValue("mode") == "1337" {
				fmt.Fprintln(w, leetify(password.Generate(p)))
			} else {
				fmt.Fprintln(w, password.Generate(p))
			}
		}
	}
}

func extractParameter(r *http.Request) password.Parameter {
	var err error
	p := password.DefaultRequest

	p.MaxLength, err = strconv.Atoi(r.FormValue("max"))
	if err != nil {
		p.MaxLength = password.DefaultMax
	}

	p.MinLength, err = strconv.Atoi(r.FormValue("min"))
	if err != nil {
		p.MinLength = password.DefaultMin
	}

	p.SpecialChars, err = strconv.Atoi(r.FormValue("spec"))
	if err != nil {
		p.SpecialChars = -1
	}

	p.Numbers, err = strconv.Atoi(r.FormValue("num"))
	if err != nil {
		p.Numbers = -1
	}

	return p
}

//leetify randomly converts a given string into leetspeak
func leetify(s string) string {
	if rand.Int()%2 == 0 {
		return s
	}

	s = strings.Replace(s, "o", "0", -1)
	s = strings.Replace(s, "O", "0", -1)

	s = strings.Replace(s, "i", "1", -1)
	s = strings.Replace(s, "I", "1", -1)

	s = strings.Replace(s, "u", "2", -1)
	s = strings.Replace(s, "U", "2", -1)

	s = strings.Replace(s, "e", "3", -1)
	s = strings.Replace(s, "E", "3", -1)

	s = strings.Replace(s, "a", "4", -1)
	s = strings.Replace(s, "A", "4", -1)

	return s
}
