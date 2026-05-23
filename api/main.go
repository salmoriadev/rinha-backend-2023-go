package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	db, err := ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository := NewPessoaRepository(db)
	handler := NewPessoaHandler(repository)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/pessoas" && r.Method == http.MethodPost {
			handler.CreatePessoa(w, r)
			return
		}

		if path == "/pessoas" && r.Method == http.MethodGet {
			handler.SearchPessoas(w, r)
			return
		}

		if strings.HasPrefix(path, "/pessoas/") && r.Method == http.MethodGet {
			handler.GetPessoa(w, r)
			return
		}

		if path == "/contagem-pessoas" && r.Method == http.MethodGet {
			handler.CountPessoas(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	log.Println("server running on port 80")

	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Fatal(err)
	}
}
