package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Book struct {
	ID        string    `json:"id" xml:"id" form:"id" query:"id"`
	Title     string    `json:"title" xml:"title" form:"title" query:"title"`
	Author    string    `json:"author" xml:"author" form:"author" query:"author"`
	Price     float64   `json:"price" xml:"price" form:"price" query:"price"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type User struct {
	ID        string    `json:"id" xml:"id" form:"id" query:"id"`
	Username  string    `json:"username" xml:"username" form:"username" query:"username"`
	Password  string    `json:"password" xml:"password" form:"password" query:"password"`
	Fullname  string    `json:"fullname" xml:"fullname" form:"fullname" query:"fullname"`
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

// use to update book
// params: id, Book
// return: json
func deleteBook(c echo.Context) error {
	id := c.Param("id")

	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var book Book
	db.Delete(&book, id)

	return c.JSON(http.StatusOK, "success")
}

func login(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	db, err := gorm.Open(mysql.Open("web_diy:Bubgum123!@/pokemon?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	result := db.First(&user, "username = ? AND password = ?", user.Username, user.Password)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func main() {
	e := echo.New()

	// Add dummy data
	// env.db.Create(&Book{Title: "Memandang Bulan", Author: "Prof. A", Price: 100})
	// env.db.Create(&Book{Title: "Memandang Bintang", Author: "Prof. B", Price: 150})

	e.POST("/login", login)

	r := e.Group("/books")

	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))

	r.GET("", getBooks)
	r.GET("/:id", getBook)
	r.POST("", createBook)
	r.PUT("/:id", updateBook)
	r.DELETE("/:id", deleteBook)

	e.Logger.Fatal(e.Start(":3000"))
}
