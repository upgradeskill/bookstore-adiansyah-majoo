package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Book struct {
	ID        string    `json:"id" xml:"id" form:"id" query:"id"`
	Title     string    `json:"title" xml:"title" form:"title" query:"title"`
	Author    string    `json:"author" xml:"author" form:"author" query:"author"`
	Price     float64   `json:"price" xml:"price" form:"price" query:"price"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

// use to get all books
// params:
// return: json
func getBooks(c echo.Context) error {
	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var books []Book
	db.Find(&books)

	return c.JSON(http.StatusOK, books)
}

// use to get specific books
// params: id string
// return: json
func getBook(c echo.Context) error {
	id := c.Param("id")

	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var book Book
	db.First(&book, id)

	return c.JSON(http.StatusOK, book)
}

// use to create book
// params: Book
// return: json
func createBook(c echo.Context) error {
	book := new(Book)
	if err := c.Bind(book); err != nil {
		return err
	}

	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Create(&book)
	return c.JSON(http.StatusCreated, "success")
}

// use to update book
// params: id, Book
// return: json
func updateBook(c echo.Context) error {
	id := c.Param("id")

	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var book Book
	bookUpdate := new(Book)
	if err := c.Bind(bookUpdate); err != nil {
		return err
	}

	// title := c.FormValue("title")
	// author := c.FormValue("author")
	// price := c.FormValue("price")

	db.First(&book, id)
	db.Model(&book).Updates(bookUpdate)

	return c.JSON(http.StatusOK, "success")
}
func main() {
	e := echo.New()

	// Add dummy data
	// env.db.Create(&Book{Title: "Memandang Bulan", Author: "Prof. A", Price: 100})
	// env.db.Create(&Book{Title: "Memandang Bintang", Author: "Prof. B", Price: 150})

	e.GET("/books", getBooks)
	e.GET("/books/:id", getBook)
	e.POST("/books", createBook)
	e.PUT("/books/:id", updateBook)

	e.Logger.Fatal(e.Start(":3000"))
}
