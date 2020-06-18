// app.go
// Package  Book API.
//
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8080
//     Version: 1.0.0
//     Contact: Vinod Pandeey<vinod.pandey1@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:meta
package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/vinodpandey1/main/model"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(book, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", book, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	fs := http.FileServer(http.Dir("./swaggerui"))
	a.Router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	// swagger:operation GET /books getBooks
	//
	// Returns all books from the system
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// responses:
	//   '200':
	//     description: book response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Book"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	a.Router.HandleFunc("/books", a.getBooks).Methods("GET")
	// swagger:operation POST /book createBook
	//
	// Creates a new book in the store.
	// Duplicates are allowed
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: book
	//   in: body
	//   description: Book to add to the store
	//   required: true
	//   schema:
	//       "$ref": "#/definitions/Book"
	//
	// responses:
	//   '200':
	//     description: book response
	//     schema:
	//       "$ref": "#/responses/bookRes"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/responses/badReq"
	a.Router.HandleFunc("/book", a.createBook).Methods("POST")
	// swagger:operation GET /book/{id} getBook
	//
	// Returns a BOOK based on a single ID,
	// if the book does not have access to the book
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// parameters:
	// - name: id
	//   in: path
	//   description: ID of book to fetch
	//   required: true
	//   type: integer
	//   format: int64
	// responses:
	//   '200':
	//     description: book response
	//     schema:
	//       "$ref": "#/definitions/Book"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/responses/badReq"
	a.Router.HandleFunc("/book/{id:[0-9]+}", a.getBook).Methods("GET")
	// swagger:operation PUT /book/{id} updateBook
	//
	// Update a  book in the store.
	// Duplicates are allowed
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: book
	//   in: body
	//   description: Book to add to the store
	//   required: true
	//   schema:
	//       "$ref": "#/definitions/Book"
	//
	// responses:
	//   '200':
	//     description: book response
	//     schema:
	//       "$ref": "#/definitions/Book"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/responses/badReq"
	a.Router.HandleFunc("/book/{id:[0-9]+}", a.updateBook).Methods("PUT")
	// swagger:operation DELETE /book/{id} deleteBook
	//
	// deletes a single Book based on the ID supplied
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   description: ID of Book to delete
	//   required: true
	//   type: integer
	//   format: int64
	// responses:
	//   '204':
	//     description: Book deleted
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/responses/badReq"
	a.Router.HandleFunc("/book/{id:[0-9]+}", a.deleteBook).Methods("DELETE")
}

func (a *App) getBooks(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := model.GetBooks(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createBook(w http.ResponseWriter, r *http.Request) {
	var u model.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.CreateBook(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	u := model.Book{ID: id}
	if err := u.GetBook(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Book not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var u model.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateBook(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Book ID")
		return
	}

	u := model.Book{ID: id}
	if err := u.DeleteBook(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
