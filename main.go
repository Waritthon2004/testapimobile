package main

import (
	"os"
	"time"

	_ "github.com/abc/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type User struct {
	Email    string
	Password string
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

var books []Book

// @title lotto
// @description This is a sample server for a lotto.
// @version 1.0
// @host https://testapimobile.onrender.com/
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	books = append(books, Book{ID: 1, Title: "NN", Author: "sdsd"})
	books = append(books, Book{ID: 2, Title: "gg", Author: "zzzz"})
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/login", loging)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Use(checkMiddleware)

	app.Get("/config", getEnv)
	app.Get("/hello", getHello)
	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)
	app.Post("/upload", uploadFile)
	app.Get("/test-html", testHtml)
	app.Listen(":8080")

}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "../upload/"+file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusNotFound).SendString("File uploads complete")
}

func testHtml(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello , world!!!",
	})
}

func checkMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func loging(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin" // example role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t, "message": "Login success"})

}

func getEnv(c *fiber.Ctx) error {
	if value, exists := os.LookupEnv("SECRET"); exists {
		return c.JSON(fiber.Map{
			"SECRET": value,
		})
	}
	return c.JSON(fiber.Map{
		"SECRET": "defalutsecret",
	})
}
