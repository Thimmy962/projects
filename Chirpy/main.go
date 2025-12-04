package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"Chirpy/internal/database"
)

type apiConfig struct{
	serverHit atomic.Int64
}


type Chirp struct {
	Body string `json:"body"`
	UserId string `json:"userId"`
}

func (chirp * Chirp) validBody() int {
	n := len(chirp.Body)
	if n > 140 {
		return 1
	}
	if n < 1 {
		return -1
	}
	return 0
}

type Server struct {
    db      *sql.DB
    queries *database.Queries
}

var profane =[]string {"kerfuffle", "sharbert", "fornax"}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	dbQuery := database.New(db)

	dbServer := &Server{
    db:      db,
    queries: dbQuery,
}

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

	serMux.HandleFunc("POST /api/validate_chirp", dbServer.CreateChirp)
	serMux.HandleFunc("POST /api/users", dbServer.createUser)
	serMux.HandleFunc("POST /admin/reset", dbServer.deleteUsers)
	serMux.HandleFunc("GET /api/chirps", dbServer.CORSMiddleware(dbServer.ListChirps))
	serMux.HandleFunc("GET /api/chirps/{id}", dbServer.GetChirp)
	serMux.HandleFunc("POST /api/login", dbServer.GetUser)


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



func (s *Server) CORSMiddleware(next func(w http.ResponseWriter, req *http.Request))  func(w http.ResponseWriter, req *http.Request) {
	return  func(w http.ResponseWriter, req *http.Request) {
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

		next(w, req)
	}
}