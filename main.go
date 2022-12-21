package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type APICep struct {
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const TIMEOUT_SECONDS = 1
const CEP_INPUT = "87711-420"

func main() {
	apiCepMsg := make(chan APICep)
	viaCepMsg := make(chan ViaCEP)

	go GetViaCEP(CEP_INPUT, viaCepMsg)
	go GetAPICep(CEP_INPUT, apiCepMsg)

	select {
	case msg := <-apiCepMsg:
		fmt.Printf("API: APICep | Message: %v", msg)
	case msg := <-viaCepMsg:
		fmt.Printf("API: ViaCEP | Message: %v", msg)
	case <-time.After(time.Second * TIMEOUT_SECONDS):
		fmt.Println("Timeout")
	}
}

func GetAPICep(cep string, apiCepMsg chan APICep) {
	resp, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	var data APICep
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	apiCepMsg <- data
}

func GetViaCEP(cep string, viaCepMsg chan ViaCEP) {
	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	viaCepMsg <- data
}
