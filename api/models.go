package main

import ()

type Pessoa struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Apelido    string   `json:"apelido"`
	Stack      []string `json:"stack"`
	Nascimento string   `json:"nascimento"`
}

type PessoaRequest struct {
	Nome       string   `json:"nome" binding:"required"`
	Apelido    string   `json:"apelido" binding:"required"`
	Stack      []string `json:"stack"`
	Nascimento string   `json:"nascimento" binding:"required"`
}

type PessoaResponse struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Apelido    string   `json:"apelido"`
	Stack      []string `json:"stack"`
	Nascimento string   `json:"nascimento"`
}
