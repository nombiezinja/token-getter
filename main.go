package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func GetCredentials() {
	cwd, _ := os.Getwd()

	err := godotenv.Load(
		filepath.Join(cwd, "./credentials.env"),
	)

	if err != nil {
		log.Fatal("Error loading env files, please configure credentials.env")
	}
}

func main() {
	hs, logger := setup()

	go func() {
		logger.Printf("Listening on http://0.0.0.0%s\n", hs.Addr)

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

		GetCredentials()
// 	https://www.facebook.com/v3.2/dialog/oauth?
//   client_id={app-id}
//   &redirect_uri={redirect-uri}
//   &state={state-param}

	request := httprutils.Request{
		Method: httprutils.Post,
		URL	"https://www.facebook.com/v3.2/dialog/oauth?client_id={app-id}&redirect_uri={redirect-uri}&state={state-param}",
		Headers: map[string]string{
			"content-Type":       "application/json",
		},
		QueryParams: map[string]string{
			"apiKey":          os.Getenv("APIKEY"),
			"verificationurl": verificationURL,
			"emailtemplate":   emailTemplate,
			"options":         options,
		},
		Body: requestBody,

	graceful(hs, logger, 5*time.Second)
}

func setup() (*http.Server, *log.Logger) {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":2019"
	}

	hs := &http.Server{Addr: addr, Handler: &server{}}

	return hs, log.New(os.Stdout, "", 0)
}

var stop = make(chan bool, 1)

func graceful(hs *http.Server, logger *log.Logger, timeout time.Duration) {

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf("Error: %v\n", err)
	} else {
		logger.Println("Server stopped")
	}
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte("Shutting down gracefully"))
	signal := true
	stop <- signal
}
