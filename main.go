package main

import "github.com/qinXpeng/rabe/rabes"

func main() {
	app := rabes.New()
	v1 := app.Group("/v1")
	{
		v1.GET("/get", func(ctx *rabes.Context) {
			a := map[string]interface{}{
				"fwvwev": "vwev",
				"vwrv":   21,
				"fwevv": map[string]interface{}{
					"pop": 213,
				},
			}
			ctx.JSON(rabes.STATUS_PAGE_OK.Code(), a)
		})
	}
	app.Run(":8081")
}
