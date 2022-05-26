package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const delay = 5
const watch = 5
const delayLogger = 100000

func main() {
	name := newSession()

	showIntroduction(name)

	for {
		showMenu()

		opt := selectOption()

		switch opt {
		case 1:
			fmt.Println("Monitorando...")
			initWatch()

		case 2:
			fmt.Println("Exibindo Logs...")
			readLogs()

		case 0:
			fmt.Println("Saindo do programa.")
			os.Exit(0)

		default:
			fmt.Println("Não conheço este comando.")
			os.Exit(-1)
		}
	}
}

func newSession() string {
	var name string
	fmt.Println("Olá senhor(a):")

	fmt.Scan(&name)

	return name
}

func showIntroduction(name string) {
	version := 1.1

	fmt.Println("Olá sr(a).", name)
	fmt.Println("Este programa está na versão", version)
}

func selectOption() int {
	var comando int
	fmt.Scan(&comando)

	fmt.Println("O comando escolhido foi:", comando)
	fmt.Println("")

	return comando
}

func showMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func initWatch() {
	sites := readSitesFile()

	for i := 0; i < watch; i++ {
		for _, site := range sites {
			testSite(site)
		}

		time.Sleep(delay * time.Second)

		fmt.Println("")
	}

	fmt.Println("")
}

func testSite(site string) {
	response, error := http.Get(site)

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}

	fmt.Println("O site", site, "está com o status", response.StatusCode)

	if response.StatusCode == 200 {
		fmt.Println("Site carregado com sucesso!")
		registerLogs(site, true, response.StatusCode)
	} else {
		fmt.Println("Site está com problemas, StatusCode:", response.StatusCode)
		registerLogs(site, false, response.StatusCode)
	}
}

func readSitesFile() []string {
	var sites []string

	file, error := os.Open("sites.txt")

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}

	reader := bufio.NewReader(file)

	for {
		linha, error := reader.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if error == io.EOF {
			break
		}
	}

	file.Close()
	return sites
}

func registerLogs(site string, status bool, code int) {
	file, error := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if error != nil {
		fmt.Println(error)
	}

	file.WriteString("[LOG] " + time.Now().Format("02/01/2006 15:04:05") + ": " + site + " - online: " + strconv.FormatBool(status) + " - status code: " + strconv.Itoa(code) + "\n")

	file.Close()
}

func readLogs() {
	file, error := os.Open("log.txt")

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}

	reader := bufio.NewReader(file)

	for {
		linha, error := reader.ReadString('\n')
		linha = strings.TrimSpace(linha)

		fmt.Println(linha)

		if error == io.EOF {
			break
		}

		time.Sleep(delayLogger * time.Microsecond)
	}
}
