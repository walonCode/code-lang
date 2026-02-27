package main

import (
	"bufio"
	"log"
	"os"

	"github.com/walonCode/code-lang/cmd/code-lang-lsp/rpc"
)

func logger(filename string) (*log.Logger){
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		panic("not a good file")
	}
	
	return log.New(logfile, "[code-lang-ls]", log.Ldate | log.Lshortfile | log.Ltime)
}

func main() {
	logger := logger("log.txt")
	logger.Println("started lsp")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Spilt)
	
	for scanner.Scan(){
		msg := scanner.Text()
		handleMessage(logger,msg )
	}
}

func handleMessage(logger *log.Logger, msg any){
	logger.Println(msg)
}
