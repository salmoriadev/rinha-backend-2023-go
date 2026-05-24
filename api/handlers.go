package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PessoaHandler struct {
	repository *PessoaRepository
}

func NewPessoaHandler(repository *PessoaRepository) *PessoaHandler {
	return &PessoaHandler{
		repository: repository,
	}
}

func (h *PessoaHandler) CreatePessoa(w http.ResponseWriter, r *http.Request) {
	req, status := ParsePessoaRequest(r)
	if status != 0 {
		w.WriteHeader(status)
		return
	}

	pessoa := Pessoa{
		ID:         uuid.NewString(),
		Apelido:    req.Apelido,
		Nome:       req.Nome,
		Nascimento: req.Nascimento,
		Stack:      req.Stack,
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	err := h.repository.InsertPessoa(ctx, pessoa)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/pessoas/"+pessoa.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *PessoaHandler) GetPessoa(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/pessoas/")

	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	pessoa, err := h.repository.GetPessoaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, pessoa)
}

func (h *PessoaHandler) SearchPessoas(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	terms, exists := values["t"]

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	term := ""
	if len(terms) > 0 {
		term = terms[0]
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	pessoas, err := h.repository.SearchPessoas(ctx, term)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, pessoas)
}

func (h *PessoaHandler) CountPessoas(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	count, err := h.repository.CountPessoas(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%d", count)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func contextWithTimeout(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 10*time.Second)
}
