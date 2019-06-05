# Tiny

个人学习 Go 而封装的简单 Web 开发库，其大部分依赖 Go 原生 http 包中的 API。

Tiny - a tiny Go web library for learning. 

## 安装（Install）

```bash
go get -u github.com/oopsguy/tiny
```

## 用法（Usage）

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/oopsguy/tiny"
)

func main() {
	server := tiny.Default()
	server.SupportTemplate("web/template/**", nil)
	server.Static("/static/", "web/static/")

	router := server.Router()
	// global-level middleware
	router.Use(func(next tiny.Handler) tiny.Handler {
		return func(ctx tiny.Ctx) error {
			fmt.Println("Global Middleware")
			return next(ctx)
		}
	})

	router.GET("", func(c tiny.Ctx) error {
		return c.Template("index.html", nil)
	}, func(next tiny.Handler) tiny.Handler {
		return func(ctx tiny.Ctx) error {
			fmt.Println("Index Middleware")
			return next(ctx)
		}
	})

	// API module
	apiRouter := router.Module("/api")
	{
		apiV1Router := apiRouter.Module("/v1")
		apiV1Router.GET("/users", func(ctx tiny.Ctx) error {
			return ctx.Content("users")
		})

		apiV2Router := apiRouter.Module("/v2")
		apiV2Router.GET("/users", func(ctx tiny.Ctx) error {
			return ctx.Content("users")
		})
	}

	// admin module
	adminRouter := router.Module("/admin")
	{
		adminRouter.Use(func(next tiny.Handler) tiny.Handler {
			return func(ctx tiny.Ctx) error {
				fmt.Println("Admin Middleware")
				return next(ctx)
			}
		})

		adminRouter.GET("/user", func(ctx tiny.Ctx) error {
			return ctx.HTML("<h1>User</h1>")
		}, func(next tiny.Handler) tiny.Handler {
			return func(ctx tiny.Ctx) error {
				fmt.Println("User Middleware")
				return next(ctx)
			}
		})

		adminRouter.GET("/posts", func(ctx tiny.Ctx) error {
			return ctx.Content("Post")
		})
	}

	if err := server.Run(":8484"); err != nil {
		log.Fatal(err)
	}
}
```

## 简介（Intro）

### 路由（Router）

```go
type Router interface {
    Handlers() []HandlerInfo
    Module(prefix string) Router
    Use(middleware ...MiddlewareHandler) Router
    Add(pattern, method string, handler Handler, middleware ...MiddlewareHandler) Router
    ANY(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    GET(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    POST(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    DELETE(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    PUT(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    PATCH(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    OPTIONS(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
    HEAD(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
}
```

### 模板渲染（Template Renderer）

默认使用 Go 的 `html.Template` 模板，你可以自己实现 `TemplateResolver` 来适配第三方模板引擎。

Go `html.Template` by default. You can register your own `TemplateResolver` to adapt third-party template engine.

```
func (t *Tiny) RegisterTemplateResolver(resolver TemplateResolver)
```   

模板渲染器接口定义与默认实现：

Interface and default implementation：

```go
type TemplateResolver interface {
	Render(name string, data interface{}) ([]byte, error)
}

type DefaultTplResolver struct {
	TemplateResolver
	Tpl *template.Template
}
```

## License

[MIT License](./LICENSE)