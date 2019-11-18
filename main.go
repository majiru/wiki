package main

import (
	"log"
	"net/http"
	"os"
	"flag"
)

const frontPage = "FrontPage"
const dataDir = "./data"

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "port to run webserver on")
	os.Mkdir(dataDir, 0755)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
