package main

import (
	"log"
	"net/http"

	"vec-node/db"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/queue-size", db.GetQueueSize).Methods("GET")
	router.HandleFunc("/queue-size", db.SetQueueSize).Methods("POST")
	router.HandleFunc("/queue-size", db.UpdateQueueSize).Methods("PUT")

	log.Printf("Server running at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}


// package main;
// import (
//     "fmt"
//     "log"
//     "net/http"
// )
// func main() {
//     http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request){
//         fmt.Fprintf(w, "Hello, World!")
//     })
//     fmt.Printf("Server running (port=8080), route: http://localhost:8080/helloworld\n")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         log.Fatal(err)
//     }
// }