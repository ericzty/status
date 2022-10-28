package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	targetsJson, err := ioutil.ReadFile("monitor.json")
	if err != nil {
		fmt.Println(err)
	}
	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	st := status{Mutex: &sync.Mutex{}}
	json.Unmarshal([]byte(targetsJson), &st.values)
	go Monitoring(st)

	engine := html.New("./views", ".html")

	server := fiber.New(fiber.Config{
		Views: engine,
	})

	server.Static("/static", "./static")

	server.Get("/", func(c *fiber.Ctx) error {
		st.Mutex.Lock()
		values := st.values
		st.Mutex.Unlock()
		location, _ := time.LoadLocation("America/New_York")
		time := time.Now().In(location).Format(time.UnixDate)
		return c.Render("index", fiber.Map{
			"Status": values,
			"Time":   time,
		})
	})

	log.Fatal(server.Listen(":3000"))
}
