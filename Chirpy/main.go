package main

import (
	"net/http"
	"log"
)




func main() {
	port := "8080"
	serMux := http.NewServeMux()

	root := http.Dir(".")
	
	serMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(root)))
	
	serMux.Handle("/app/assets", http.StripPrefix("/app/assets", http.FileServer(root+"/assets")))
	serMux.HandleFunc("/healthz", healthz)

	server := http.Server{Addr: ":"+ port, Handler: serMux}

	log.Printf("Serving files on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


func healthz(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("status-code", "200 OK")

	header.Write(w)
	
	

	w.Write([]byte("OK"))
}