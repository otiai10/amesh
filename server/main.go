package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ResponseJSON(w, map[string]interface{}{
			"message": "Hello, Amesh Server",
		})
	})
	err := http.ListenAndServe(":4010", nil)
	log.Fatalln(err)
}

func ResponseJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.Write(b)
}
