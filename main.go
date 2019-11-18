package main

import (
	"log"
	"net/http"
	"os"
	"flag"

	"github.com/majiru/aitm"
)

const frontPage = "FrontPage"
const dataDir = "./data"

var port string
var auth string

func init() {
	flag.StringVar(&port, "port", "8080", "port to run webserver on")
	flag.StringVar(&auth, "auth", "", "config file for aitm")
	os.Mkdir(dataDir, 0755)
}

func main() {
	flag.Parse()
	mux := &http.ServeMux{}
	mux.HandleFunc("/", rootHandler)
	srv := &http.Server{
		Addr: ":"+port,
		Handler: mux,
	}
	if auth == "" {
		log.Fatal(srv.ListenAndServe())
	}
	asrv := aitm.WrapServer(srv)
	f, err := os.Open(auth)
	if err != nil {
		log.Fatalf("could not open %s: %v\n", auth, err)
	}
	err = asrv.LoadUsers(f)
	if err != nil {
		log.Fatal("err parsing users conf:", err)
	}
	f.Close()
	log.Fatal(asrv.ListenAndServe())
}
