package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PessoaRepository struct {
	db *pgxpool.Pool
}

func NewPessoaRepository(db *pgxpool.Pool) *PessoaRepository {
	return &PessoaRepository{
		db: db,
	}
}

func (r *PessoaRepository) InsertPessoa(ctx context.Context, pessoa Pessoa) error {
	var stackJSON any = nil

	if pessoa.Stack != nil {
		data, err := json.Marshal(pessoa.Stack)
		if err != nil {
			log.Printf("Error marshaling stack: %v", err)
			return err
		}

		stackJSON = string(data)
	}

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO pessoas
			(id, nome, apelido, stack, nascimento, busca)
		VALUES
			($1, $2, $3, $4::jsonb, $5, $6)
		`,
		pessoa.ID,
		pessoa.Nome,
		pessoa.Apelido,
		stackJSON,
		pessoa.Nascimento,
		buildBusca(pessoa),
	)

	if err != nil {
		log.Printf("Error inserting pessoa: %v", err)
		return err
	}

	return nil
}

func (r *PessoaRepository) GetPessoaByID(ctx context.Context, id string) (PessoaResponse, error) {
	var pessoa PessoaResponse
	var stackRaw []byte

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id::text,
			nome,
			apelido,
			stack,
			to_char(nascimento, 'YYYY-MM-DD')
		FROM pessoas
		WHERE id = $1
		`,
		id,
	).Scan(
		&pessoa.ID,
		&pessoa.Nome,
		&pessoa.Apelido,
		&stackRaw,
		&pessoa.Nascimento,
	)

	if err != nil {
		log.Printf("Error fetching pessoa by ID: %v", err)
		return PessoaResponse{}, err
	}

	if stackRaw != nil {
		err = json.Unmarshal(stackRaw, &pessoa.Stack)
		if err != nil {
			return PessoaResponse{}, err
		}
	}

	return pessoa, nil
}

func (r *PessoaRepository) SearchPessoas(ctx context.Context, term string) ([]PessoaResponse, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id::text,
			nome,
			apelido,
			stack,
			to_char(nascimento, 'YYYY-MM-DD')
		FROM pessoas
		WHERE busca LIKE '%' || lower($1) || '%'
		LIMIT 50
		`,
		term,
	)

	if err != nil {
		log.Printf("Error searching pessoas: %v", err)
		return nil, err
	}

	defer rows.Close()

	pessoas := make([]PessoaResponse, 0)

	for rows.Next() {
		var pessoa PessoaResponse
		var stackRaw []byte

		err := rows.Scan(
			&pessoa.ID,
			&pessoa.Nome,
			&pessoa.Apelido,
			&stackRaw,
			&pessoa.Nascimento,
		)

		if err != nil {
			log.Printf("Error scanning pessoa row: %v", err)
			return nil, err
		}

		if stackRaw != nil {
			err = json.Unmarshal(stackRaw, &pessoa.Stack)
			if err != nil {
				log.Printf("Error unmarshaling stack for pessoa ID %s: %v", pessoa.ID, err)
				return nil, err
			}
		}

		pessoas = append(pessoas, pessoa)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over pessoa rows: %v", err)
		return nil, err
	}

	return pessoas, nil
}

func (r *PessoaRepository) CountPessoas(ctx context.Context) (int, error) {
	var count int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM pessoas
		`,
	).Scan(&count)

	if err != nil {
		log.Printf("Error counting pessoas: %v", err)
		return 0, err
	}

	return count, nil
}
