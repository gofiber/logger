// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasttemplate"
)

// Filter variables
const (
	strTime          = "time"
	strReferer       = "referer"
	strProtocol      = "protocol"
	strIp            = "ip"
	strIps           = "ips"
	strHost          = "host"
	strMethod        = "method"
	strPath          = "path"
	strUrl           = "url"
	strUa            = "ua"
	strLatency       = "latency"
	strStatus        = "status"
	strBody          = "body"
	strBytesSent     = "bytesSent"
	strBytesReceived = "bytesReceived"
	strRoute         = "route"
	strError         = "error"
	strHeader        = "header:"
	strQuery         = "query:"
	strForm          = "form:"
	strCookie        = "cookie:"
)

// Config ...
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool
	// Format defines the logging format with defined variables
	// Optional. Default: "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n"
	// Possible values:
	// time, ip, ips, url, host, method, path, protocol, route
	// referer, ua, latency, status, body, error, bytesSent, bytesReceived
	// header:<key>, query:<key>, form:<key>, cookie:<key>
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
	var mutex sync.RWMutex
	// Middleware settings
	tmpl := fasttemplate.New(cfg.Format, "${", "}")
	timestamp := time.Now().Format(cfg.TimeFormat)
	// Update date/time every second in a seperate go routine
	if strings.Contains(cfg.Format, "${time}") {
		go func() {
			for {
				mutex.Lock()
				timestamp = time.Now().Format(cfg.TimeFormat)
				mutex.Unlock()
				time.Sleep(250 * time.Millisecond)
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
		// Get new buffer
		buf := bytebufferpool.Get()
		_, err := tmpl.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case strTime:
				mutex.RLock()
				defer mutex.RUnlock()
				return buf.WriteString(timestamp)
			case strReferer:
				return buf.WriteString(c.Get(fiber.HeaderReferer))
			case strProtocol:
				return buf.WriteString(c.Protocol())
			case strIp:
				return buf.WriteString(c.IP())
			case strIps:
				return buf.WriteString(c.Get(fiber.HeaderXForwardedFor))
			case strHost:
				return buf.WriteString(c.Hostname())
			case strMethod:
				return buf.WriteString(c.Method())
			case strPath:
				return buf.WriteString(c.Path())
			case strUrl:
				return buf.WriteString(c.OriginalURL())
			case strUa:
				return buf.WriteString(c.Get(fiber.HeaderUserAgent))
			case strLatency:
				return buf.WriteString(stop.Sub(start).String())
			case strStatus:
				return buf.WriteString(strconv.Itoa(c.Fasthttp.Response.StatusCode()))
			case strBody:
				return buf.WriteString(c.Body())
			case strBytesReceived:
				return buf.WriteString(strconv.Itoa(len(c.Fasthttp.Request.Body())))
			case strBytesSent:
				return buf.WriteString(strconv.Itoa(len(c.Fasthttp.Response.Body())))
			case strRoute:
				return buf.WriteString(c.Route().Path)
			case strError:
				return buf.WriteString(c.Error().Error())
			default:
				switch {
				case strings.HasPrefix(tag, strHeader):
					return buf.WriteString(c.Get(tag[7:]))
				case strings.HasPrefix(tag, strQuery):
					return buf.WriteString(c.Query(tag[6:]))
				case strings.HasPrefix(tag, strForm):
					return buf.WriteString(c.FormValue(tag[5:]))
				case strings.HasPrefix(tag, strCookie):
					return buf.WriteString(c.Cookies(tag[7:]))
				}
			}
			return 0, nil
		})
		if err != nil {
			_, _ = buf.WriteString(err.Error())
		}
		if _, err := cfg.Output.Write(buf.Bytes()); err != nil {
			fmt.Println(err)
		}
		bytebufferpool.Put(buf)
	}
}
