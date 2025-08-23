package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type InputData struct {
	Data string `json:"data"`
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func cmpRegEx(find string, regex string) bool {
	temp := regexp.MustCompile(regex)

	if temp.MatchString(find) {
		return true
	}
	return false
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "main.html", gin.H{
			"title": "Intext Testing",
		})

		r.Static("/src", "./src")
	})

	r.POST("/api/send", func(c *gin.Context) {
		var input InputData

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		code := input.Data
		var output []byte
		os.WriteFile("temp.itx", []byte(code), 0644)

		file, err := os.Open("temp.itx")
		Check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if cmpRegEx(line, `(read|write|append)\(.+\);?`) {
				output = []byte(fmt.Sprintf("You cannot use file i/o funcs! -> %v", scanner.Text()))
			}
		}

		erra := scanner.Err()
		Check(erra)

		cmd := exec.Command("./ITX-CLI", "run", "temp.itx")
		var errb error
		if string(output) == "" {
			output, errb = cmd.Output()
			Check(errb)
		}

		c.JSON(http.StatusOK, gin.H{"message": string(output)})
	})

	r.Run(":8080")
}
