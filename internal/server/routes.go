package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.HandleFunc("/*", s.handleRedisProxy)

	return r
}

func (s *Server) handleRedisProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		s.sendError(w, "No command specified", http.StatusBadRequest)
		return
	}

	cmd := strings.ToUpper(parts[0])
	args := parts[1:]

	if r.Body != nil {
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			var bodyArgs []string
			if err := json.Unmarshal(body, &bodyArgs); err == nil {
				args = append(args, bodyArgs...)
			}
		}
	}

	ctx := context.Background()
	result, err := s.executeCommand(ctx, cmd, args)
	if err != nil {
		s.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendSuccess(w, result)
}

func (s *Server) executeCommand(ctx context.Context, cmd string, args []string) (interface{}, error) {
	redisArgs := make([]interface{}, len(args)+1)
	redisArgs[0] = cmd
	for i, arg := range args {
		redisArgs[i+1] = arg
	}
	result := s.rdb.Do(ctx, redisArgs...)
	return result.Result()
}

func (s *Server) sendSuccess(w http.ResponseWriter, result interface{}) {
	json.NewEncoder(w).Encode(Response{Result: result})
}

func (s *Server) sendError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Error: msg})
}
