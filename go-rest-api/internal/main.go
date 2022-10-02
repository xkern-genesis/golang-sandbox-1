package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/xkern-genesis/golang-sandbox-1.git/go-rest-api/pkg/swagger/server/restapi"
	"github.com/xkern-genesis/golang-sandbox-1.git/go-rest-api/pkg/swagger/server/restapi/operations"
)

const TERMINAL_NAME = "[xkern-gen]$ "

func main() {
	/**
	// Defining main api router
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%sHello,%q\n", TERMINAL_NAME, html.EscapeString(r.URL.Path))
	})
	//Log running api
	log.Println("Listenin on   localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	**/

	//Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	// Final process
	defer func() {
		if err := server.Shutdown(); err != nil {
			// Handling error
			log.Fatalln(err)
		}
	}()

	// Server features
	server.Port = 8080
	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(GetHealthHandler)
	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUserHandler)
	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherHandler)
	// Starting the server
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

/**
	Health handler function
**/
func GetHealthHandler(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

/**
	GetHello User returns Hello + your_name
**/
func GetHelloUserHandler(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}

/**
	Get gopher in png
**/
func GetGopherHandler(gopher operations.GetGopherNameParams) middleware.Responder {
	var URL string
	if gopher.Name != "" {
		URL = "https://github.com/scraly/gophers/raw/main/" + gopher.Name + ".png"
	} else {
		// by default we return dr wh gopher
		URL = "https://github.com/scraly/gophers/raw/main/dr-who.png"
	}
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}
