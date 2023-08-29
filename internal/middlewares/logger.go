package middlewares

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	colors := map[string]*color.Color{
		"TRACE": color.New(color.FgHiWhite),
		"DEBUG": color.New(color.FgHiGreen),
		"INFO":  color.New(color.FgHiBlue),
		"WARN":  color.New(color.FgHiYellow),
		"ERROR": color.New(color.FgHiRed),
	}

	methodColors := map[string]*color.Color{
		"GET":    color.New(color.FgHiCyan),
		"POST":   color.New(color.FgHiMagenta),
		"PUT":    color.New(color.FgHiYellow),
		"DELETE": color.New(color.FgHiRed),
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logLevel := "INFO" // Set the default log level

		// Determine the log level based on status code
		if statusCode >= 500 {
			logLevel = "ERROR"
		} else if statusCode >= 400 {
			logLevel = "WARN"
		}

		colorizer, found := colors[logLevel]
		if !found {
			colorizer = color.New() // Default color for unknown log levels
		}

		methodColorizer, methodFound := methodColors[method]
		if !methodFound {
			methodColorizer = color.New() // Default color for unknown methods
		}

		fmt.Printf("[GIN] %v |%3d| %s | %s | %-7s %s | %s\n",
			time.Now().Format(time.RFC1123),
			statusCode,
			colorizer.SprintfFunc()(logLevel),
			methodColorizer.SprintfFunc()(method),
			c.ClientIP(),
			path,
			latency,
		)
	}
}
