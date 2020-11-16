package main

import (
	"log"
	"net"
	"net/http"

	"github.com/JovidYnwa/http/cmd/app"
	"github.com/JovidYnwa/http/pkg/banners"
	"github.com/JovidYnwa/http/pkg/server"
)

func main() {
	/* 	banner := []banners.Banner{{ID: 1,
	   		Titile:  "somet title",
	   		Content: "some content",
	   		Button:  "some button",
	   		Link:    "some link",
	   	}}
	   	log.Println(banner) */

	mux := http.NewServeMux()
	greetingSvc := banners.NewService()
	server := app.NewServer(mux, greetingSvc)
	server.Init()

	httpServer := http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: server, //myHandler
	}

	log.Println("server start")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("http serv error: ", err)
	}
}

/* log.Println("handler{}.ServerHttp()")
writer.Write([]byte("hello http"))
*/

//lecture 22
func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))

	srv.Register("/payments/{id}/{ds}", func(req *server.Request) {
		id := req.PathParams["id"]
		log.Print(id)
	})
	log.Print("server run in ", host+":"+port)
	return srv.Start()
}
