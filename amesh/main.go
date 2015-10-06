// main
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/otiai10/amesh"
	"github.com/otiai10/query"
)

var (
	port  = "4010"
	serve = false
)

func init() {
	flag.BoolVar(&serve, "s", false, "run the server")
	flag.StringVar(&port, "p", port, "port of server")
	flag.Parse()
}

func main() {
	if serve {
		server()
	}
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ResponseJSON(w, amesh.GetEntry(), r)
	})
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatalln(err)
}

// ResponseJSON ...
func ResponseJSON(w http.ResponseWriter, v interface{}, r *http.Request) {
	var b []byte
	var err error
	if query.Bool(r, "pretty", false) {
		b, err = json.MarshalIndent(v, "", "\t")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.Write(b)
}
