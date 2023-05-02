package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	host := "db" //os.Getenv("DB_HOST")
	user := "root" //os.Getenv("DB_USER")
	password := "root_password_vec" //os.Getenv("DB_PASSWORD")
	database := "app_db" //os.Getenv("DB_NAME")
	port := "3306" //os.Getenv("DB_PORT")

	var err error
	db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+database+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// Set up the queue size table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS queue_size (id INT AUTO_INCREMENT PRIMARY KEY, size INT NOT NULL DEFAULT 0)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetQueueSize(w http.ResponseWriter, r *http.Request) {
	var size int
	err := db.QueryRow("SELECT size FROM queue_size WHERE id = 1").Scan(&size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + strconv.Itoa(size) + `}`))
}

func SetQueueSize(w http.ResponseWriter, r *http.Request) {

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request couldn't decode", http.StatusBadRequest)
		return
	}

	size, ok := data["size"]
	if !ok {
		http.Error(w, "Invalid Request no 'size'", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO queue_size (id, size) VALUES (1, ?) ON DUPLICATE KEY UPDATE size = ?", size, size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + size + `}`))
}


// func SetQueueSize(w http.ResponseWriter, r *http.Request) {
// 	var size int
// 	if err := json.NewDecoder(r.Body).Decode(&size); err != nil {
// 		log.Println(err)
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	_, err := db.Exec("INSERT INTO queue_size (id, size) VALUES (1, ?) ON DUPLICATE KEY UPDATE size = ?", size, size)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"size":` + strconv.Itoa(size) + `}`))
// }

func UpdateQueueSize(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request couldn't decode", http.StatusBadRequest)
		return
	}
	log.Println("hello 12")
	size, ok := data["size"]
	if !ok {
		http.Error(w, "Invalid Request no 'size'", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE queue_size SET size = ? WHERE id = 1", size)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"size":` + size + `}`))
}
