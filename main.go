package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// Converter número no formato AAAAMMDD para string no formato DD-MM-AAAA
func converterData(numero string) time.Time {
	// Converter a string para um objeto time.Time
	data, err := time.Parse("20060102", numero)
	if err != nil {
		panic(err)
	}

	// Formatar a data no formato desejado
	return data
}

type Registro struct {
	Sistema   string    `json:"sistema"`
	Empresa   string    `json:"empresa"`
	Tabela    string    `json:"tabela"`
	Data      time.Time `json:"data"`
	Ano       string    `json:"ano"`
	Mes       string    `json:"mes"`
	Dia       string    `json:"dia"`
	Registros int       `json:"registros"`
}

func carregarDados() []Registro {
	file, err := excelize.OpenFile("./volumetria/volumetria.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rows, err := file.GetRows("LGM")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var registros []Registro

	for _, row := range rows {
		var registro Registro
		for i, colCell := range row {
			if i == 0 {
				registro.Sistema = colCell
			} else if i == 1 {
				registro.Empresa = colCell
			} else if i == 2 {
				registro.Tabela = colCell
			} else if _, err := strconv.Atoi(colCell); err == nil && len(colCell) == 8 {
				registro.Data = converterData(colCell) // Assume que qualquer coluna com 8 dígitos é uma data
				ano, mes, dia := extrairAnoMesDia(colCell)
				registro.Ano = ano
				registro.Mes = mes
				registro.Dia = dia
			} else if i == 4 { // Supondo que a quinta coluna contém o valor dos registros
				fmt.Printf("Data: %s\n", colCell)
				registro.Registros, err = strconv.Atoi(colCell)
				if err != nil {
					fmt.Printf("Erro ao converter registro para inteiro: %v\n", err)
					return nil
				}
			}
		}
		registros = append(registros, registro)
	}

	return registros
}

func extrairAnoMesDia(numero string) (string, string, string) {

	// Extrair o ano e o mês
	ano := numero[:4]  // Retorna "2023"
	mes := numero[4:6] // Retorna "01"
	dia := numero[6:8] // Retorna "01"

	return ano, mes, dia
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data := carregarDados()
	if data == nil {
		http.Error(w, "Erro ao carregar dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Servidor iniciado em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
