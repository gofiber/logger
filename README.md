### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/logger
```
### Example
```go
package main

import 
  "github.com/gofiber/fiber"
  "github.com/gofiber/logger"
)

func main() {
  app := fiber.New()

  app.Use(logger.New())

  app.Get("/", func(c *fiber.Ctx) {
    c.Send("Welcome!")
  })

  app.Listen(3000)
}
```
### Test
```curl
curl http://localhost:3000
```
