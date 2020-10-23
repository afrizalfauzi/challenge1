package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// model
type Kontak struct {
	ID           string `json:"Id"`
	NamaDepan    string `json:"NamaDepan"`
	NamaBelakang string `json:"NamaBelakang"`
	NoHp         string `json:"NoHp"`
	Email        string `json:"Email"`
	Alamat       string `json:"Alamat"`
	Umur         string `json:"Umur"`
}

func getKontak(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var kontaks []Kontak

	sql := `SELECT
				Id,
				IFNULL(NamaDepan,''),
				IFNULL(NamaBelakang,''),
				IFNULL(NoHp,''),
				IFNULL(Email,''),
				IFNULL(Alamat,''),
				IFNULL(Umur,'')
			FROM kontak`
	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var kontak Kontak
		err := result.Scan(&kontak.ID, &kontak.NamaDepan, &kontak.NamaBelakang, &kontak.NoHp, &kontak.Email, &kontak.Alamat, &kontak.Umur)

		if err != nil {
			panic(err.Error())
		}
		kontaks = append(kontaks, kontak)
	}
	json.NewEncoder(w).Encode(kontaks)
}

func createKontak(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		ID := r.FormValue("Id")
		NamaDepan := r.FormValue("NamaDepan")
		NamaBelakang := r.FormValue("NamaBelakang")

		stmt, err := db.Prepare("INSERT INTO kontak (Id,NamaDepan, NamaBelakang) VALUES (?,?,?)")

		_, err = stmt.Exec(ID, NamaDepan, NamaBelakang)

		if err != nil {
			fmt.Fprint(w, "Data Ganda")
		} else {
			fmt.Fprint(w, "Data Dibuat")
		}
	}
}

func getKontaks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var kontaks []Kontak
	params := mux.Vars(r)

	sql := `SELECT
				Id,
				IFNULL(NamaDepan,''),
				IFNULL(NamaBelakang,''),
				IFNULL(NoHp,''),
				IFNULL(Email,''),
				IFNULL(Alamat,''),
				IFNULL(Umur,'')
			FROM kontak WHERE Id = ?`
	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var kontak Kontak

	for result.Next() {

		err := result.Scan(&kontak.ID, &kontak.NamaDepan, &kontak.NamaBelakang, &kontak.NoHp, &kontak.Email, &kontak.Alamat, &kontak.Umur)

		if err != nil {
			panic(err.Error())
		}

		kontaks = append(kontaks, kontak)
	}
	json.NewEncoder(w).Encode(kontaks)

}

func updateKontak(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNamaDepan := r.FormValue("NamaDepan")

		stmt, err := db.Prepare("UPDATE kontak SET NamaDepan = ? WHERE Id = ?")

		_, err = stmt.Exec(newNamaDepan, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data tidak ditemukan or Request error")
		}

		fmt.Fprintf(w, "kontak with Id = %s was updated", params["id"])
	}
}

func deleteKontak(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM kontak WHERE Id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "Gagal menghapus")
	}

	fmt.Fprintf(w, "Kontak with Id = %s was deleted", params["id"])
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var kontaks []Kontak

	ID := r.FormValue("Id")
	NamaDepan := r.FormValue("NamaDepan")

	sql := `SELECT
	Id,
	IFNULL(NamaDepan,''),
	IFNULL(NamaBelakang,''),
	IFNULL(NoHp,''),
	IFNULL(Email,''),
	IFNULL(Alamat,''),
	IFNULL(Umur,'')
	FROM kontak WHERE Id = ? AND NamaDepan = ?`

	result, err := db.Query(sql, ID, NamaDepan)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var kontak Kontak

	for result.Next() {
		err := result.Scan(&kontak.ID, &kontak.NamaDepan, &kontak.NamaBelakang, &kontak.NoHp, &kontak.Email, &kontak.Alamat, &kontak.Umur)

		if err != nil {
			panic(err.Error())
		}

		kontaks = append(kontaks, kontak)
	}

	json.NewEncoder(w).Encode(kontaks)

}

func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/challenge")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/kontak", getKontak).Methods("GET")
	r.HandleFunc("/kontak/{id}", getKontak).Methods("GET")
	r.HandleFunc("/kontak", createKontak).Methods("POST")
	r.HandleFunc("/kontak/{id}", updateKontak).Methods("PUT")
	r.HandleFunc("/kontak/{id}", deleteKontak).Methods("DELETE")

	r.HandleFunc("/getkontak", getPost).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))
}
