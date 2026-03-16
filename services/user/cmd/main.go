package main

import "net/http"

func main() {

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OKOKOK"))
		if err != nil {
			return
		}
	}
	http.HandleFunc("/users", handler)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		return
	}

}
