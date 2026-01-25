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
	secret string
}


type Chirp struct {
	Body string `json:"body"`
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
	secret string
	apiKey string
}

var profane =[]string {"kerfuffle", "sharbert", "fornax"}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

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
	secret: os.Getenv("JWT_TOKEN_SECRET"),
}

	port := "8888"
	serMux := http.NewServeMux()

	server := http.Server{Addr: ":"+ port, Handler: serMux}
	cfg := apiConfig{secret: os.Getenv("JWT_TOKEN_SECRET")}

	root := http.Dir(".")
	
	serMux.Handle("/api/app/", cfg.middlewareMetricsInc(http.StripPrefix("/api/app/", http.FileServer(root))))
	serMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(root))))
	serMux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(root+"/assets"))))

	serMux.HandleFunc("GET /api/healthz", healthz)
	serMux.HandleFunc("GET /healthz", healthz)

	serMux.Handle("GET /admin/metrics", metrics(&cfg))

	serMux.HandleFunc("POST /api/users", dbServer.CORSMiddleware(dbServer.createUser))
	serMux.HandleFunc("POST /admin/reset", dbServer.CORSMiddleware(dbServer.deleteUsers))
	serMux.HandleFunc("GET /api/chirps", dbServer.CORSMiddleware(dbServer.listChirps))
	serMux.HandleFunc("POST /api/chirps", dbServer.CORSMiddleware(dbServer.createChirp))
	serMux.HandleFunc("GET /api/chirps/{id}", dbServer.CORSMiddleware(dbServer.getChirp))
	serMux.HandleFunc("POST /api/login", dbServer.CORSMiddleware(dbServer.getUserToken))
	serMux.HandleFunc("GET /api/getuser", dbServer.CORSMiddleware(dbServer.getUserDet))
	serMux.HandleFunc("POST /api/refresh", dbServer.CORSMiddleware(dbServer.refresh))
	serMux.HandleFunc("PUT /api/users", dbServer.CORSMiddleware(dbServer.editUserDetail))
	serMux.HandleFunc("DELETE /api/chirps/{id}", dbServer.CORSMiddleware(dbServer.deleteChirp))
	serMux.HandleFunc("POST /api/polka/webhooks", dbServer.CORSMiddleware(dbServer.webhook))


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


func (s *Server) CORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

		next(w, req)
	}
}
