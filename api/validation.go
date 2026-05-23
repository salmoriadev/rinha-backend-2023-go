package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func ParsePessoaRequest(r *http.Request) (PessoaRequest, int) {
	var raw map[string]json.RawMessage

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return PessoaRequest{}, http.StatusBadRequest
	}
	if len(body) == 0 {
		return PessoaRequest{}, http.StatusBadRequest
	}

	err = json.Unmarshal(body, &raw)
	if err != nil {
		return PessoaRequest{}, http.StatusBadRequest
	}
	nome, status := parseRequiredString(raw, "nome")
	if status != 0 {
		return PessoaRequest{}, status
	}
	apelido, status := parseRequiredString(raw, "apelido")
	if status != 0 {
		return PessoaRequest{}, status
	}
	nascimento, status := parseRequiredString(raw, "nascimento")
	if status != 0 {
		return PessoaRequest{}, status
	}
	stack, status := parseStack(raw)
	if status != 0 {
		return PessoaRequest{}, status
	}

	req := PessoaRequest{
		Nome:       nome,
		Apelido:    apelido,
		Nascimento: nascimento,
		Stack:      stack,
	}
	if status := validatePessoaRequest(req); status != 0 {
		return PessoaRequest{}, status
	}
	return req, 0
}

func parseRequiredString(raw map[string]json.RawMessage, field string) (string, int) {
	valueRaw, exists := raw[field]
	if !exists {
		return "", http.StatusUnprocessableEntity
	}
	value := string(valueRaw)
	if value == "null" || value == "" {
		return "", http.StatusUnprocessableEntity
	}

	var result string

	if err := json.Unmarshal(valueRaw, &result); err != nil {
		return "", http.StatusBadRequest
	}
	return result, 0
}

func parseStack(raw map[string]json.RawMessage) ([]string, int) {
	valueRaw, exists := raw["stack"]
	if !exists {
		return nil, 0
	}
	if string(valueRaw) == "null" {
		return nil, 0
	}

	var rawItens []json.RawMessage
	if err := json.Unmarshal(valueRaw, &rawItens); err != nil {
		return nil, http.StatusBadRequest
	}
	stack := make([]string, 0, len(rawItens))
	for _, item := range rawItens {
		if string(item) == "null" {
			return nil, http.StatusBadRequest
		}
		var stackItem string
		if err := json.Unmarshal(item, &stackItem); err != nil {
			return nil, http.StatusBadRequest
		}
		if utf8.RuneCountInString(stackItem) > 32 {
			return nil, http.StatusUnprocessableEntity
		}
		if stackItem == "" {
			return nil, http.StatusUnprocessableEntity
		}

		stack = append(stack, stackItem)
	}
	return stack, 0
}

func validatePessoaRequest(req PessoaRequest) int {
	if req.Nome == "" || req.Apelido == "" || req.Nascimento == "" {
		return http.StatusUnprocessableEntity
	}

	if utf8.RuneCountInString(req.Nome) > 100 {
		return http.StatusUnprocessableEntity
	}
	if utf8.RuneCountInString(req.Apelido) > 32 {
		return http.StatusUnprocessableEntity
	}

	if _, err := time.Parse("2006-01-02", req.Nascimento); err != nil {
		return http.StatusUnprocessableEntity
	}

	for _, item := range req.Stack {
		if item == "" {
			return http.StatusUnprocessableEntity
		}
		if utf8.RuneCountInString(item) > 32 {
			return http.StatusUnprocessableEntity
		}
	}

	return 0
}

func buildBusca(pessoa Pessoa) string {
	parts := []string{
		pessoa.Apelido,
		pessoa.Nome,
	}

	if pessoa.Stack != nil {
		parts = append(parts, pessoa.Stack...)
	}

	return strings.ToLower(strings.Join(parts, " "))
}
