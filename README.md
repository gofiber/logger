# Logger

![Release](https://img.shields.io/github/release/gofiber/logger.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/gofiber/logger/workflows/Test/badge.svg)
![Security](https://github.com/gofiber/logger/workflows/Security/badge.svg)
![Linter](https://github.com/gofiber/logger/workflows/Linter/badge.svg)

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/logger
```
### Format
`Format` defines the logging format with defined variables
Default: "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n"  

Possible values: time, ip, ips, url, host, method, path, protocol, route, referer, ua, latency, status, body, error, bytesSent, bytesReceived, header:<key>, query:<key>, form:<key>, cookie:<key>

### Example
```go
package main

import (
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
