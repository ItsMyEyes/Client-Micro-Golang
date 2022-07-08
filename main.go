package main

import (
	pb "client/common/model"
	"client/resthandlers"
	"client/routes"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	port     int
	authAddr string
)

func init() {
	flag.IntVar(&port, "port", 9003, "api service port")
	flag.StringVar(&authAddr, "auth_addr", ":3000", "authenticaton service address")
	flag.Parse()
}

func main() {

	conn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	authSvcClient := pb.NewAuthServiceClient(conn)
	authHandlers := resthandlers.NewAuthHandlers(authSvcClient)
	authRoutes := routes.NewAuthRoutes(authHandlers)

	router := mux.NewRouter().StrictSlash(true)
	routes.Install(router, authRoutes)

	log.Printf("API service running on 0.0.0.0:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), routes.WithCORS(router)))
}
