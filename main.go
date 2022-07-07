package main

import (
	"fmt"
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
// params: echo.Context
// return: json
func getBooks(c echo.Context) error {
	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var books []Book
	db.Find(&books)
	fmt.Println("{}", books)

	return c.JSON(http.StatusOK, books)
}

func main() {
	e := echo.New()

	// Add dummy data
	// env.db.Create(&Book{Title: "Memandang Bulan", Author: "Prof. A", Price: 100})
	// env.db.Create(&Book{Title: "Memandang Bintang", Author: "Prof. B", Price: 150})

	e.GET("/books", getBooks)

	e.Logger.Fatal(e.Start(":3000"))
}
