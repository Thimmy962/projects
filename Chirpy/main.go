package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	serverHit atomic.Int64
}


type Chirp struct {
	Body string `json:"body"`
}

type Error struct {
	Err string `json:"error"`
}


var profane =[]string {"kerfuffle", "sharbert", "fornax"}

func main() {
	port := "8888"
	serMux := http.NewServeMux()

	server := http.Server{Addr: ":"+ port, Handler: serMux}
	cfg := apiConfig{}

	root := http.Dir(".")
	
	serMux.Handle("/api/app/", cfg.middlewareMetricsInc(http.StripPrefix("/api/app/", http.FileServer(root))))
	serMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(root))))
	serMux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(root+"/assets"))))

	serMux.HandleFunc("GET /api/healthz", healthz)

	serMux.Handle("GET /admin/metrics", metrics(&cfg))
	serMux.HandleFunc("POST /admin/reset", cfg.reset())

	serMux.HandleFunc("POST /api/validate_chirp", ValidateChirp)


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


func metrics(cfg *apiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "text/html")
		data := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"
        hits := fmt.Sprintf(data, cfg.serverHit.Load())
        w.Write([]byte(hits))
    }
}


//middleware
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.serverHit.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) reset() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.serverHit.Store(0)
		// w.Header().Set("status-code")
		metrics(cfg)
	})
}


func ValidateChirp(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	n, err := buf.ReadFrom(req.Body)
	// if there is an error in the writing
	if err != nil {
		Err_500ApplicationJson(w, "something went wrong")
		return
	}
	// if the number of bytes read(n) is > 140 
	if n > 140 {
		Err_400ApplicationJson(w, "chirp too long")
		return
	}

	chirp := Chirp{Body: buf.String()}
	jsonData, err := json.Marshal(chirp)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	buf.Reset()
	// convert jsonData to string from byte, pass it to a function that returns a string
	// convert the returned string to byte and svave in jsonData 
	jsonData = []byte(profaneFUnc(string(jsonData)))
	buf.Write(jsonData)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(buf.Bytes())
}


