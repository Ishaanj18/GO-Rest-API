package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	apipath = "/apis/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string // these wont be getting exported or getting public
}
type Book struct {
	Id, Name, Isbn string // these are camel cased as we will be exporting them or making them public
}

func main() {
	//env variables
	dbHost := os.Getenv("DB_Host")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbName := os.Getenv("DB_Name")
	if dbName == "" {
		dbName = "Library"
	}
	dbPass := os.Getenv("DB_Pass")
	if dbPass == "" {
		dbPass = "jameswatt@18"
	}
	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = apipath
	}
	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods("GET")
	r.HandleFunc(apiPath, l.postBooks).Methods("POST")
	http.ListenAndServe(":8080", r)
}
func (l library) postBooks(w http.ResponseWriter, r *http.Request) {
	// retrieve the request
	books := Book{}
	json.NewDecoder(r.Body).Decode(&books)
	//open connection
	db := l.openConnection()
	// write the data
	insertQuery, err := db.Prepare("insert into books values (?,?,?)")
	if err != nil {
		log.Fatalf("preparing the query resulted in the following error %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while beginning the transaction got the following error %s \n", err.Error())
	}

	_, err = tx.Stmt(insertQuery).Exec(books.Id, books.Name, books.Isbn)
	if err != nil {
		log.Fatalf("error while executing the query %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("error while commiting the transaction %s\n", err.Error())
	}
	//close the connection
	l.closeConnection(db)

}
func (l library) getBooks(w http.ResponseWriter, r *http.Request) { // all the functions openconnection , get books are method of library struct and that is why have (l library)
	//open connection
	db := l.openConnection()
	//read the books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the books table :%s \n", err.Error())
	}

	books := []Book{}
	// rows.next() will basically be true until the rows get finished. basically a while loop of numnber of rows < 0 and then number of rows--.
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn) // storing rows values to the respective variables we created
		if err != nil {
			log.Fatalf("while scanning the row %s\n", err.Error())
		}
		// instance of book struct
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}

	json.NewEncoder(w).Encode(books)

	//close the connection
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("opening the connection to the database %s \n", err.Error()) // err.Error() converts error to string
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s \n", err.Error())
	}
}
