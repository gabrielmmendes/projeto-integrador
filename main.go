package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"strings"
)

var db = fazConexaoComBanco()
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	// Configuração do servidor para servir arquivos estáticos (HTML, CSS, JS, imagens, etc.)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/paciente", pacientes)

	alimentaBancoDeDados()

	log.Println("Server rodando na porta 8080")
	// Inicia o servidor na porta 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func pacientes (w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		return
	}
	busca := strings.TrimSpace(r.Form.Get("busca"))

	Pacientes := buscaPacientePorNome(busca)

	templates.ExecuteTemplate(w, "pacientes.html", Pacientes)
}

func fazConexaoComBanco() *sql.DB {
	// carrega arquivo .env com dados de conexão com o banco
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env")
	}

	// faz a busca dos atributos no arquivo .env para usa-las na conexão com banco
	usuarioBancoDeDados := os.Getenv("USUARIO")
	senhaDoUsuario := os.Getenv("SENHA")
	nomeDoBancoDeDados := os.Getenv("NOME_BANCO_DE_DADOS")
	dadosParaConexao := "user=" + usuarioBancoDeDados + " dbname=" + nomeDoBancoDeDados + " password=" + senhaDoUsuario + " host=localhost port=5432 sslmode=disable"
	database, err := sql.Open("postgres", dadosParaConexao)
	if err != nil {
		fmt.Println(err)
	}

	// cria tabela paciente com atributos como: id, nome, cpf, data de nascimento, telefone, sexo e booleanos referente a situação fisica
	_, err = database.Query("CREATE TABLE IF NOT EXISTS paciente (id SERIAL PRIMARY KEY, nome VARCHAR(255) UNIQUE NOT NULL, cpf VARCHAR(15) UNIQUE NOT NULL, data_nascimento VARCHAR(12), telefone_celular VARCHAR(20), sexo VARCHAR(10), esta_fumante boolean, faz_uso_alcool boolean, esta_situacao_rua boolean)")
	if err != nil {
		log.Fatal(err)
	}

	return database
}

func cadastraPaciente(paciente Paciente) {
	// insere paciente no banco de dados
	_, err := db.Exec(`INSERT INTO paciente (nome, cpf, data_nascimento, telefone_celular, sexo, esta_fumante, faz_uso_alcool, esta_situacao_rua) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) on conflict do nothing`, paciente.Nome, paciente.Cpf, paciente.DataNascimento, paciente.Telefone, paciente.Sexo, paciente.EstaFumante, paciente.FazUsoAlcool, paciente.EstaSituacaoDeRua)
	if err != nil {
		fmt.Println(err)
	}
}

func buscaPacientePorNome(nome string) Pacientes {
	// retorna pacientes por nome
	busca, err := db.Query(`SELECT * FROM paciente WHERE nome LIKE concat('%', text($1), '%')`, nome)
	if err != nil {
		fmt.Println(err)
	}

	var pacientes Pacientes

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for busca.Next() {

		var paciente Paciente

		// Armazena os valores em variáveis
		var Id uint64
		var Nome, Cpf, DataNascimento, Telefone, Sexo string
		var EstaFumante, FazUsoAlcool, EstaSituacaoDeRua bool

		// Faz o Scan do SELECT
		err = busca.Scan(&Id, &Nome, &Cpf, &DataNascimento, &Telefone, &Sexo, &EstaFumante, &FazUsoAlcool, &EstaSituacaoDeRua)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		paciente.Id = Id
		paciente.Nome = Nome
		paciente.Cpf = Cpf
		paciente.DataNascimento = DataNascimento
		paciente.Telefone = Telefone
		paciente.Sexo = Sexo
		paciente.EstaFumante = EstaFumante
		paciente.FazUsoAlcool = FazUsoAlcool
		paciente.EstaSituacaoDeRua = EstaSituacaoDeRua

		// Junta a Struct com Array
		pacientes.Pacientes = append(pacientes.Pacientes, paciente)
	}

	return pacientes
}

func alimentaBancoDeDados() {
	var Pacientes Pacientes

	// lê o arquivo paciente.json e passa para o Go
	jsonFile, _ := os.Open("paciente.json")
	byteJson, _ := io.ReadAll(jsonFile)

	err := json.Unmarshal(byteJson, &Pacientes)
	if err != nil {
		log.Fatal(err)
	}

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
