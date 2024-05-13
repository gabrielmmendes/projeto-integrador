package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
)

var db = fazConexaoComBanco()

func main() {
	// Configuração do servidor para servir arquivos estáticos (HTML, CSS, JS, imagens, etc.)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	alimentaBancoDeDados()

	log.Println("Server rodando na porta 8080")
	// Iniciar o servidor na porta 8080
	http.ListenAndServe(":8080", nil)
}

func fazConexaoComBanco() *sql.DB {
	dadosParaConexao := "user=postgres dbname=postgres password=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dadosParaConexao)
	if err != nil {
		fmt.Println(err)
	}

	db.Query("CREATE TABLE IF NOT EXISTS paciente (id SERIAL PRIMARY KEY, nome VARCHAR(255) UNIQUE NOT NULL, cpf VARCHAR(15) UNIQUE NOT NULL, data_nascimento VARCHAR(12), telefone_celular VARCHAR(20), sexo VARCHAR(10), esta_fumante boolean, faz_uso_alcool boolean, esta_situacao_rua boolean)")

	return db
}

func cadastraPaciente(paciente Paciente) {
	_, err := db.Exec(`INSERT INTO paciente (nome, cpf, data_nascimento, telefone_celular, sexo, esta_fumante, faz_uso_alcool, esta_situacao_rua) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, paciente.Nome, paciente.Cpf, paciente.DataNascimento, paciente.Telefone, paciente.Sexo, paciente.EstaFumante, paciente.FazUsoAlcool, paciente.EstaSituacaoDeRua)
	if err != nil {
		fmt.Println(err)
	}
}

func alimentaBancoDeDados() {
	var Pacientes Pacientes

	jsonFile, _ := os.Open("paciente.json")
	byteJson, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteJson, &Pacientes)

	for i := 0; i < len(Pacientes.Pacientes); i++ {
		cadastraPaciente(Pacientes.Pacientes[i])
	}
}

type Paciente struct {
	Id                uint64
	Nome              string `json:"nome"`
	Cpf               string `json:"cpf"`
	DataNascimento    string `json:"data_nasc"`
	Telefone          string `json:"celular"`
	Sexo              string `json:"sexo"`
	EstaFumante       bool   `json:"esta_fumante"`
	FazUsoAlcool      bool   `json:"faz_uso_alcool"`
	EstaSituacaoDeRua bool   `json:"esta_situacao_de_rua"`
}

type Pacientes struct {
	Pacientes []Paciente `json:"pacientes"`
}
