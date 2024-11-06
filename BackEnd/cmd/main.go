package main

import (
	"fmt"
	"net/http"

	"github.com/agrawalpratham/Connectability/BackEnd/config"
	"github.com/agrawalpratham/Connectability/BackEnd/database"
	"github.com/agrawalpratham/Connectability/BackEnd/internal"
	"github.com/gorilla/sessions"
)

func main() {
	//Assigning the port number
	config.App.PortNumber = ":8080"


	//Setting the server configuration
	svr := &http.Server{
		Addr: config.App.PortNumber,
		Handler: internal.Routes(),
	}

	//Making connection to database
	d, err := database.ConnectSQL("host=localhost port=5432 dbname=ConnectabilityDatabase user=postgres password=pratham")
	if err != nil{
		fmt.Println(fmt.Sprintf("Error connecting to database : %v", err))
	}else{
		fmt.Println("Successfully connected to database")
		config.App.DBConnection = d
	}
	defer d.SQL.Close()

	//Setting the session configuration
	session := sessions.NewCookieStore([]byte("6tGbPGST2hUprIHea5xzyTmXAdw7xpfI"))
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteLaxMode
	config.App.Session = session

	//Starting the server
	fmt.Println("Starting the server on port 8080")
	err = svr.ListenAndServe()
	if( err != nil){
		fmt.Println("Error starting the server : ", err)
	}
}