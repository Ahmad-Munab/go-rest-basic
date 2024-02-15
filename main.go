package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Paradozical Sajid", Author: "Arif Azad", Quantity: 4},
	{ID: "2", Title: "Paradozical Sajid", Author: "Arif Azad", Quantity: 2},
	{ID: "3", Title: "Jemon Torun Chai", Author: "Masud Rana", Quantity: 3},
}

func getBooks(req *gin.Context) {
	req.IndentedJSON(http.StatusOK, books)
	return
}

func getBookById(req *gin.Context) {
	id := req.Param("id")
	book, err := bookById(id)

	if err != nil {
		req.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	req.IndentedJSON(http.StatusOK, book)
	return
}

func bookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("Book not found")
}

func deleteBookById(id string) ([]book, error) {
	var remainingBooks []book
	found := false
	for i, b := range books {
		if b.ID == id {
			found = true
			remainingBooks = append(remainingBooks, books[i+1:]...)
			fmt.Println(remainingBooks)
			break
		}
		remainingBooks = append(remainingBooks, b)
	}
	if found == false {
		return nil, errors.New("Could not find the book")
	}

	books = remainingBooks
	return books, nil
}

func createBooks(req *gin.Context) {
	var newBooks []book

	if err := req.BindJSON(&newBooks); err != nil {
		req.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Arguiments"})
		return
	}

	books = append(books, newBooks...)
	req.IndentedJSON(http.StatusCreated, books)
	return
}

func checkoutBooks(req *gin.Context) {
	var books_ids []string

	if err := req.BindJSON(&books_ids); err != nil {
		req.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Arguiments"})
		return
	}

	var purchasedBooks []*book
	for _, id := range books_ids {
		book, err := bookById(id)

		if err != nil {
			req.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		if book.Quantity <= 0 {
			req.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book " + id + " is not available"})
			return
		}
		purchasedBooks = append(purchasedBooks, book)
	}

	for _, book := range purchasedBooks {
		book.Quantity -= 1
	}

	req.IndentedJSON(http.StatusOK, purchasedBooks)
	return
}

func returnBooks(req *gin.Context) {
	var books_ids []string

	if err := req.BindJSON(&books_ids); err != nil {
		req.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Arguiments"})
		return
	}

	var returnedBooks []*book
	for _, id := range books_ids {
		book, err := bookById(id)

		if err != nil {
			req.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		returnedBooks = append(returnedBooks, book)
	}

	for _, book := range returnedBooks {
		book.Quantity += 1
	}

	req.IndentedJSON(http.StatusOK, gin.H{"message": "Books Added"})
	return
}

func removeBook(req *gin.Context) {
	var book_id string

	if err := req.BindJSON(&book_id); err != nil {
		req.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Arguiments"})
		return
	}

	remainingBooks, err := deleteBookById(book_id)
	if err != nil {
		req.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	req.IndentedJSON(http.StatusOK, remainingBooks)
	return
}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", createBooks)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkout", checkoutBooks)
	router.PATCH("/return", returnBooks)
	router.DELETE("/books", removeBook)

	router.Run("0.0.0.0:3004")
}
