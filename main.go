// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber"
	"github.com/valyala/fasttemplate"
)

// Config ...
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool
	// Format defines the logging format with defined variables
	// Optional. Default: "${time} - ${ip} - ${method} ${path}\t${ua}\n"
	// Possible values: time, ip, url, host, method, path, protocol
	// referer, ua, header:<key>, query:<key>, formform:<key>, cookie:<key>
	Format string
	// TimeFormat https://programming.guide/go/format-parse-string-time-date-example.html
	// Optional. Default: 15:04:05
	TimeFormat string
	// Output is a writter where logs are written
	// Default: os.Stderr
	Output io.Writer
}

// New ...
func New(config ...Config) func(*fiber.Ctx) {
	// Init config
	var cfg Config
	// Set config if provided
	if len(config) > 0 {
		cfg = config[0]
	}
	// Set config default values
	if cfg.Format == "" {
		cfg.Format = "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n"
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = "15:04:05"
	}
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	// Middleware settings
	tmpl := fasttemplate.New(cfg.Format, "${", "}")
	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}
	timestamp := time.Now().Format(cfg.TimeFormat)
	// Update date/time every second in a seperate go routine
	if strings.Contains(cfg.Format, "${time}") {
		go func() {
			for {
				timestamp = time.Now().Format(cfg.TimeFormat)
				time.Sleep(1 * time.Second)
			}
		}()
	}
	// Middleware function
	return func(c *fiber.Ctx) {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			c.Next()
			return
		}
		start := time.Now()
		// handle request
		c.Next()
		// build log
		stop := time.Now()
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		defer pool.Put(buf)
		_, err := tmpl.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time":
				return buf.WriteString(timestamp)
			case "referer":
				return buf.WriteString(c.Get(fiber.HeaderReferer))
			case "protocol":
				return buf.WriteString(c.Protocol())
			case "ip":
				return buf.WriteString(c.IP())
			case "host":
				return buf.WriteString(c.Hostname())
			case "method":
				return buf.WriteString(c.Method())
			case "path":
				return buf.WriteString(c.Path())
			case "url":
				return buf.WriteString(c.OriginalURL())
			case "ua":
				return buf.WriteString(c.Get(fiber.HeaderUserAgent))
			case "latency":
				return buf.WriteString(stop.Sub(start).String())
			case "status":
				return buf.WriteString(strconv.Itoa(c.Fasthttp.Response.StatusCode()))
			default:
				switch {
				case strings.HasPrefix(tag, "header:"):
					return buf.WriteString(c.Get(tag[7:]))
				case strings.HasPrefix(tag, "query:"):
					return buf.WriteString(c.Query(tag[6:]))
				case strings.HasPrefix(tag, "form:"):
					return buf.WriteString(c.FormValue(tag[5:]))
				case strings.HasPrefix(tag, "cookie:"):
					return buf.WriteString(c.Cookies(tag[7:]))
				}
			}
			return 0, nil
		})
		if err != nil {
			buf.WriteString(err.Error())
		}
		if _, err := cfg.Output.Write(buf.Bytes()); err != nil {
			fmt.Println(err)
		}
	}
}
