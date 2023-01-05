package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type book struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

var books []book = []book{
	book{ID: "1", Title: "The Alchemist", Author: "Paulo Coelho", Quantity: 10},
	book{ID: "2", Title: "The Monk Who Sold His Ferrari", Author: "Robin Sharma", Quantity: 5},
	book{ID: "3", Title: "The Secret", Author: "Rhonnda Byrne", Quantity: 7},
}

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Welcome to the Book Store",
	})
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookById(id string) (*book, error) {
	for _, book := range books {
		if book.ID == id {
			return &book, nil
		}
	}

	err := errors.New("Book not found")
	return nil, err
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var newBook book
	currBook, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	currBook.Title = newBook.Title
	currBook.Author = newBook.Author
	currBook.Quantity = newBook.Quantity
	c.IndentedJSON(http.StatusOK, currBook)
}

func checkoutBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Book is out of stock",
		})
		return
	}
	
	book.Quantity -= 1
	for i, b := range books {
		if b.ID == id {
			books[i] = *book
		}
	}
	c.IndentedJSON(http.StatusOK, book)

}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func notFound(c *gin.Context) {
	err := errors.New("Not Found")
	c.JSON(http.StatusNotFound, gin.H{
		"message": err.Error(),
	})
}


func main() {
	r := gin.Default()
	r.GET("/", home)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/books", getBooks)
	r.GET("/books/:id", getBook)
	r.POST("/books", createBook)
	r.PUT("/books/:id", updateBook)
	r.PATCH("/books/:id", checkoutBook)

	// r.GET("*", notFound)

	r.Run(":3000")
}