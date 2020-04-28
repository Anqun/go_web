package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	// For each matched request Context will hold the route definition
	//router.POST("/user/:name/*action", func(c *gin.Context) {
	//	c.FullPath() == "/user/:name/*action" // true
	//})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.POST("/form_post", func(c *gin.Context) {

		//Content-Type: application/json 方式在body中获取参数
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("read body err, %v\n", err)
			return
		}
		println("json:", string(body))
		// Content-Type: application/x-www-form-urlencoded 方式在表单中获取参数
		//POST /post?id=1234&page=1 HTTP/1.1
		//Content-Type: application/x-www-form-urlencoded
		//name=manu&message=this_is_great
		//nick := c.Request.PostFormValue("nick")
		//message := c.Request.PostFormValue("message")、

		//POST /post?ids[a]=1234&ids[b]=hello HTTP/1.1
		//Content-Type: application/x-www-form-urlencoded
		//names[first]=thinkerou&names[second]=tianou
		//ids := c.QueryMap("ids")
		//names := c.PostFormMap("names")
		var testBody TestBody
		if err = json.Unmarshal(body, &testBody); err != nil {
			fmt.Printf("Unmarshal err, %v\n", err)
			return
		}

		nick := testBody.Nick
		message := testBody.Message
		println(fmt.Sprintf("%d", testBody))

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
			"nick":    nick,
			"body":    testBody,
		})
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.MaxMultipartMemory = 1 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/uploads", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	//Grouping routes
	router.Run(":8080")
}

//定义类时,如果需要json转对象，那么变量名称需要大写，否则转换失败
type TestBody struct {
	Nick    string
	Message string
}
