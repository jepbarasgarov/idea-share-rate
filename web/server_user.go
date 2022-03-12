package web

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func (s *Server) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
// 	handleName := "HandleUserDelete"

// 	ctx := r.Context()
// 	ipAddress, err := helpers.GetIP(r)
// 	clog := log.WithContext(ctx).WithFields(log.Fields{
// 		"remote-addr": ipAddress,
// 		"uri":         r.RequestURI,
// 	})

// 	requestLang := helpers.GetRequestLang(r)

// 	if err != nil {
// 		eMsg := "couldn't get ip address"
// 		clog.WithError(err).Error(eMsg)
// 		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
// 		errs.SendResponse(w, err, nil, clog, requestLang)
// 		return
// 	}
// 	roles := []responses.UserRole{
// 		responses.UserRoleAdmin,
// 	}

// 	cu, err := s.UserRequirments(ctx, w, r, roles)
// 	if err != nil {
// 		eMsg := "UserRequirments error in " + handleName
// 		clog.WithError(err).Error(eMsg)
// 		errs.SendResponse(w, err, nil, clog, requestLang)
// 		return
// 	}

// 	id, err := uuid.FromString(mux.Vars(r)["id"])
// 	if err != nil {
// 		eMsg := "user ID is not compatible"
// 		clog.WithError(err).Error(eMsg)
// 		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
// 		errs.SendResponse(w, err, nil, clog, cu.Language)
// 		return
// 	}
// 	err = s.c.UserDelete(ctx, cu, id.String())
// 	if err != nil {
// 		eMsg := "error in s.c.UserDelete"
// 		clog.WithError(err).Error(eMsg)
// 		errs.SendResponse(w, err, nil, clog, cu.Language)
// 		return
// 	}
// 	err = responses.ErrOK
// 	errs.SendResponse(w, err, nil, clog, cu.Language)
// 	clog.Info(handleName + " success")
// }

/// GET

func (s *Server) HandleUserAutocomleteList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserAutocomleteList"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
		responses.UserRoleUser,
	}

	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.UserAutocompleteList(ctx, cu)
	if err != nil {
		eMsg := "error in s.c.UserList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	Resp := make([]responses.UserLightData, 0)
	for _, us := range *data {
		respUser := responses.UserLightData{
			ID:        us.ID,
			Role:      us.Role,
			Firstname: us.Firstname,
			Lastname:  us.Lastname,
		}
		Resp = append(Resp, respUser)
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, cu.Language)
	clog.Info(handleName + " success")
}

//////////////////////////////////MONGO/////////////////////////////////////////////////////////////////////

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserLogin"

	ctx := r.Context()
	requestLang := helpers.GetRequestLang(r)

	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	m := make(map[string]string)
	m["username"] = username
	m["password"] = password

	names, err1 := helpers.VerifyMinLen(m)
	if err1 != nil {
		eMsg := "len == 0 -- > " + names
		clog.WithError(err1).Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	data, tokens, err := s.c.UserLogin(ctx, username, password)
	if err != nil {
		eMsg := "error in s.c.UserLogin"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Response := responses.UserLogin{
		ID:           data.ID,
		Username:     data.Username,
		Firstname:    data.Firstname,
		Lastname:     data.Lastname,
		Role:         data.Role,
		Status:       responses.Active,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Response, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleUserGiveToken(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserGiveToken"

	ctx := r.Context()

	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	refreshToken := r.FormValue("refresh_token")
	token := r.Header.Get("Authorization")

	if !strings.HasPrefix(token, "Bearer ") {
		clog.Error("User is not logged in")
		err := errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		errs.SendResponse(w, err, nil, clog, requestLang)
	}

	username, _ := helpers.VerifyAccessToken(token[7:], config.Conf.Jwt.Secret)

	data, err := s.c.UserGiveAccessToken(ctx, username, refreshToken)
	if err != nil {
		eMsg := "error in s.c.UserGiveAccessToken"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Response := responses.Tokens{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Response, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserCreate"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
	}
	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	var user models.UserCreate

	user.Username = r.FormValue("username")
	if len(user.Username) == 0 || len(user.Username) > 32 {
		eMsg := "username length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}
	password := r.FormValue("password")
	if !helpers.IsPasswordValid(password) {
		eMsg := "possword is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}
	user.Firstname = r.FormValue("firstname")
	if len(user.Firstname) == 0 || len(user.Firstname) > 32 {
		eMsg := "firstname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Lastname = r.FormValue("lastname")
	if len(user.Lastname) == 0 || len(user.Lastname) > 32 {
		eMsg := "lastname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Role, err = helpers.ConvertStringToUserRole(r.FormValue("role"))
	if err != nil {
		eMsg := "error in helpers.ConvertStringToUserRole(role)"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.UserCreate(ctx, &user, password, cu)
	if err != nil {
		eMsg := "error in s.c.UserCreate()"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserUpdate"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
	}
	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	var user models.UserUpdate

	user.ID, err = primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		eMsg := "invalid user ID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Username = r.FormValue("username")
	if len(user.Username) == 0 || len(user.Username) > 32 {
		eMsg := "username length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Firstname = r.FormValue("firstname")
	if len(user.Firstname) == 0 || len(user.Firstname) > 32 {
		eMsg := "firstname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Lastname = r.FormValue("lastname")
	if len(user.Lastname) == 0 || len(user.Lastname) > 32 {
		eMsg := "lastname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Status, err = helpers.ConvertStringToUserStatus(r.FormValue("status"))
	if err != nil {
		eMsg := "error in helpers.ConvertStringToUserStatus"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	user.Role, err = helpers.ConvertStringToUserRole(r.FormValue("role"))
	if err != nil {
		eMsg := "error in helpers.ConvertStringToUserRole(role)"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}
	if user.ID == cu.ID {
		eMsg := "You are not allowed to change your own data here"
		clog.Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_FB_owndata_USER)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.UserUpdate(ctx, &user)
	if err != nil {
		eMsg := "error in s.c.UserUpdate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleUserUpdateOwnPassword(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserUpdateOwnPassword"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
		responses.UserRoleUser,
	}
	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")

	if !helpers.IsPasswordValid(newPassword) {
		emsg := "password is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = s.c.UserUpdateOwnPassword(ctx, cu, oldPassword, newPassword)
	if err != nil {
		eMsg := "error in s.c.UserUpdateOwnPassword"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, cu.Language)
	clog.Info(handleName + " success")
}

func (s *Server) HandleAdminUpdatePassword(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleAdminUpdatePassword"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
	}

	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		eMsg := "userID is not compatible"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	newPassword := r.FormValue("password")
	if !helpers.IsPasswordValid(newPassword) {
		eMsg := "password is not compatible"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = s.c.AdminUpdatePassword(ctx, cu, id, newPassword)
	if err != nil {
		eMsg := "error in s.c.AdminUpdatePassword"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, cu.Language)
	clog.Info(handleName + " success")
}

// GET

func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleUserGet"

	ctx := r.Context()
	ipAddress, err := helpers.GetIP(r)
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"remote-addr": ipAddress,
		"uri":         r.RequestURI,
	})

	requestLang := helpers.GetRequestLang(r)

	if err != nil {
		eMsg := "couldn't get ip address"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	roles := []responses.UserRole{
		responses.UserRoleAdmin,
		responses.UserRoleUser,
	}

	cu, err := s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	ID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		eMsg := "user ID is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.UserGet(ctx, cu, ID)
	if err != nil {
		eMsg := "error in s.c.UserList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, cu.Language)
	clog.Info(handleName + " success")
}
