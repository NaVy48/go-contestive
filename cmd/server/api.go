package main

import (
	"contestive/api/handler"
	apimiddleware "contestive/api/middleware"
	"contestive/api/payload"
	"contestive/config"
	"contestive/repository/postgresql"
	"contestive/service/auth"
	"contestive/service/contest"
	"contestive/service/judgemanager"
	"contestive/service/jwt"
	"contestive/service/password"
	"contestive/service/problem"
	"contestive/service/submission"
	"contestive/service/user"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	router chi.Router
	logger *log.Logger
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// HandleAPI returns an API handler
func HandleAPI(cfg *config.Config) http.Handler {

	s := server{chi.NewRouter(), log.Default()}
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))

	store := newStore(cfg)
	// judgeManager := judge.NewJudgeManager(judge.ManagerConfiguration{
	// 	Address: cfg.Judge.Address,
	// 	Secrets: cfg.Judge.JudgeCredentials,
	// 	Repo:    store,
	// })

	if cfg.JWT.Expiration == 0 {
		cfg.JWT.Expiration = 5 * 60 * 60 // 5 hours
	}
	jwtService, err := jwt.NewJWTService(cfg.JWT.Secret, time.Duration(cfg.JWT.Expiration)*time.Second)
	if err != nil {
		s.logger.Fatalln(err)
	}

	jsonHandler := payload.NewJSONHandler(s.logger)
	passwordService := password.NewService()
	userService := user.NewService(store.UserRepository(), passwordService)
	problemService := problem.NewService(store.ProblemRepository())
	contestService := contest.NewService(store.ContestRepository())
	authService := auth.NewService(userService, jwtService, passwordService)
	jm := judgemanager.NewJudgeManager(cfg.Judge.Address, cfg.Judge.JudgeCredentials, store.SubmissionRepository(), store.ProblemRepository())
	submissionService := submission.NewService(store.SubmissionRepository(), jm, problemService, contestService)

	s.router.Mount("/auth", handler.NewAuthHandler(jsonHandler, authService))

	s.router.Group(func(r chi.Router) {
		r.Use(apimiddleware.Auth(jsonHandler, authService))
		r.Mount("/users", handler.NewUserHandler(jsonHandler, userService))
		r.Mount("/problems", handler.NewProblemHandler(jsonHandler, problemService))
		r.Mount("/contests", handler.NewContestHandler(jsonHandler, contestService))
		r.Mount("/submissions", handler.NewSubmissionHandler(jsonHandler, submissionService))
	})

	return &s
}

func newStore(cfg *config.Config) postgresql.Store {
	config := cfg.Database.PostgreSQL
	store, err := postgresql.Connect(config.Address, config.Username, config.Password, config.Database)
	if err != nil {
		panic(err.Error())
	}
	return store
}
