package web

import (
	"belli/onki-game-ideas-mongo-backend/api"
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	c *api.APIController
}

func NewServer(apiController *api.APIController) *Server {
	return &Server{
		c: apiController,
	}
}

func NewRouter(s *Server) *mux.Router {
	r := mux.NewRouter()

	//_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-

	// User //
	r.HandleFunc("/api/v1/user/login", s.HandleUserLogin).Methods("POST")
	r.HandleFunc("/api/v1/user/token", s.HandleUserGiveToken).Methods("POST")
	r.HandleFunc("/api/v1/user/create", s.HandleUserCreate).Methods("POST")
	r.HandleFunc("/api/v1/user/update", s.HandleUserUpdate).Methods("POST")
	r.HandleFunc("/api/v1/user/own/change-password", s.HandleUserUpdateOwnPassword).Methods("POST")
	r.HandleFunc("/api/v1/user/change-password/{id}", s.HandleAdminUpdatePassword).Methods("POST")
	r.HandleFunc("/api/v1/user/get/{id}", s.HandleUserGet).Methods("POST")
	r.HandleFunc("/api/v1/user/autocomplete/list", s.HandleUserAutocomleteList).Methods("POST")
	r.HandleFunc("/api/v1/user/delete/{id}", s.HandleUserDelete).Methods("POST")

	// WORKER
	r.HandleFunc("/api/v1/worker/autocomplete/list", s.HandleWorkerAutocompleteList).Methods("POST")
	r.HandleFunc("/api/v1/worker/position/list", s.HandlePositionList).Methods("POST")
	r.HandleFunc("/api/v1/worker/create", s.HandleWorkerCreate).Methods("POST")
	r.HandleFunc("/api/v1/worker/update", s.HandleWorkerUpdate).Methods("POST")
	r.HandleFunc("/api/v1/worker/delete/{id}", s.HandleWorkerDelete).Methods("POST")
	r.HandleFunc("/api/v1/worker/get/{id}", s.HandleWorkerGetByID).Methods("POST")

	//IDEA
	r.HandleFunc("/api/v1/idea/create", s.HandleIdeaCreate).Methods("POST")
	r.HandleFunc("/api/v1/idea/get/{id}", s.HandleIdeaGet).Methods("POST")
	r.HandleFunc("/api/v1/idea/delete/{id}", s.HandleIdeaDelete).Methods("POST")
	r.HandleFunc("/api/v1/idea/update/{id}", s.HandleIdeaUpdate).Methods("POST")
	r.HandleFunc("/api/v1/idea/list", s.HandleIdeaList).Methods("POST")
	r.HandleFunc("/api/v1/idea/rate", s.HandleRateIdea).Methods("POST")

	// r.HandleFunc("/api/v1/idea/get/{id}/pdf", s.HandleIdeaGetPdf).Methods("POST")
	r.HandleFunc("/api/v1/idea/list/pdf", s.HandleIdeaListGetPdf).Methods("POST")

	//CRITERIA
	r.HandleFunc("/api/v1/idea/criteria/create", s.HandleCriteriaCreate).Methods("POST")
	r.HandleFunc("/api/v1/idea/criteria/update", s.HandleCriteriaUpdate).Methods("POST")
	r.HandleFunc("/api/v1/idea/criteria/delete/{id}", s.HandleCriteriaDelete).Methods("POST")
	r.HandleFunc("/api/v1/idea/criteria/list", s.HandleCriteriaList).Methods("POST")

	//IDEA-ADDITIONALS
	r.HandleFunc("/api/v1/idea/genre/list", s.HandleGenreList).Methods("POST")
	r.HandleFunc("/api/v1/idea/genre/create", s.HandleGenreCreate).Methods("POST")
	r.HandleFunc("/api/v1/idea/genre/update", s.HandleGenreUpdate).Methods("POST")
	r.HandleFunc("/api/v1/idea/genre/delete", s.HandleGenreDelete).Methods("POST")

	r.HandleFunc("/api/v1/idea/mechanic/create", s.HandleMechanicsCreate).Methods("POST")
	r.HandleFunc("/api/v1/idea/mechanic/update", s.HandleMechanicUpdate).Methods("POST")
	r.HandleFunc("/api/v1/idea/mechanic/delete", s.HandleMechanicDelete).Methods("POST")
	r.HandleFunc("/api/v1/idea/mechanic/list", s.HandleMechanicList).Methods("POST")

	r.HandleFunc("/api/v1/datetime-sync", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("X-Date", time.Now().UTC().String())
	}).Methods("HEAD")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}

func (s *Server) UserRequirments(ctx context.Context, w http.ResponseWriter, r *http.Request,
	roles []responses.UserRole) (user *responses.ActionInfo, err error) {
	clog := log.WithFields(log.Fields{
		"handle": "UserRequirments",
	}).WithContext(ctx)

	token := r.Header.Get("Authorization")

	if !strings.HasPrefix(token, "Bearer ") {
		clog.Error("User is not logged in")
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return nil, err
	}

	username, err := helpers.VerifyAccessToken(token[7:], config.Conf.Jwt.Secret)
	if err != nil {
		clog.Error("User is not logged in")
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return nil, err
	}

	us, err := s.c.UserGetByUsername(ctx, username)
	if err != nil {
		eMsg := "error in s.c.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	if us == nil || us.Status == responses.Blocked {
		clog.Error("User is not logged in")
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return nil, err
	}

	ok := false
	for _, role := range roles {
		if us.Role == role {
			ok = true
			break
		}
	}

	if !ok {
		clog.Error("user role doesn't fit")
		err = errs.NewHttpErrorForbidden(errs.ERR_FB_USER)
		return nil, err
	}
	user = &responses.ActionInfo{}
	lang := r.Header.Get("X-Lang")
	switch lang {
	case "tk":
		user.Language = responses.Turkmen
	case "ru":
		user.Language = responses.Russian
	default:
		user.Language = responses.English
	}

	user.ID = us.ID
	user.Username = us.Username
	user.Firstname = us.Firstname
	user.Lastname = us.Lastname
	user.Role = us.Role
	user.Status = us.Status

	return
}
