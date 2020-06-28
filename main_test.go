// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/valyala/bytebufferpool"
)

func TestNew_withRoutePath(t *testing.T) {
	routePath := "/test/:param/sufix"
	format := "route=${route}"
	expectedOutput := "route=/test/:param/sufix"

	// fake output
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	n := New(Config{
		Format: format,
		Output: buf,
	})
	app := fiber.New(&fiber.Settings{DisableStartupMessage: true})
	app.Use(n)

	app.Get(routePath, func(ctx *fiber.Ctx) {
		ctx.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/af593469-3133-4943-b193-31f02e6e82e9/sufix", nil)

	_, err := app.Test(req, 1000)
	if err != nil {
		t.Errorf("Has: %+v, expected: nil", err)
	}

	if buf.String() != expectedOutput {
		t.Errorf("Has: %s, expected: %s", buf.String(), expectedOutput)
	}
}
