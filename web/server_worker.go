package web

import (
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

/////////////////////////////////////////////////////////////////////////////MONGO///////////////////////////////////////////////////

func (s *Server) HandleWorkerCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleWorkerCreate"

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

	var worker models.WorkerCreate

	worker.Firstname = r.FormValue("firstname")
	if len(worker.Firstname) == 0 || len(worker.Firstname) > 64 {
		eMsg := "firstname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	worker.Lastname = r.FormValue("lastname")
	if len(worker.Lastname) == 0 || len(worker.Lastname) > 64 {
		eMsg := "Lastname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	worker.Position = r.FormValue("position")
	if len(worker.Position) == 0 || len(worker.Position) > 128 {
		eMsg := "Position length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.WorkerCreate(ctx, &worker, cu)
	if err != nil {
		eMsg := "error in s.c.WorkerCreate()"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	Resp := responses.WorkerLightData{
		ID:        data.ID,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Position:  data.Position,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleWorkerGetByID(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleWorkerGetByID"

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

	id := mux.Vars(r)["id"]

	data, err := s.c.WorkerGetByID(ctx, id, cu)
	if err != nil {
		eMsg := "error in s.c.WorkerGetByID"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Resp := responses.WorkerLightData{
		ID:        data.ID,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Position:  data.Position,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleWorkerDelete(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleWorkerDelete"

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

	ID := mux.Vars(r)["id"]

	err = s.c.WorkerDelete(ctx, ID, cu)
	if err != nil {
		eMsg := "error in s.c.WorkerDelete()"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleWorkerUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleWorkerUpdate"

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

	var worker models.WorkerUpdate

	worker.ID = r.FormValue("id")

	worker.Firstname = r.FormValue("firstname")
	if len(worker.Firstname) == 0 || len(worker.Firstname) > 64 {
		eMsg := "firstname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	worker.Lastname = r.FormValue("lastname")
	if len(worker.Lastname) == 0 || len(worker.Lastname) > 64 {
		eMsg := "Lastname length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}
	worker.Position = r.FormValue("position")
	if len(worker.Position) == 0 || len(worker.Position) > 128 {
		eMsg := "Position length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.WorkerUpdate(ctx, &worker, cu)
	if err != nil {
		eMsg := "error in s.c.WorkerUpdate()"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	Resp := responses.WorkerLightData{
		ID:        data.ID,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Position:  data.Position,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleWorkerAutocompleteList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleWorkerAutocompleteList"

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

	data, err := s.c.WorkerAutoCompleteList(ctx)
	if err != nil {
		eMsg := "error in s.c.WorkerAutoCompleteList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	Resp := make([]responses.WorkerLightData, 0)
	for _, worker := range *data {
		respWorker := responses.WorkerLightData{
			ID:        worker.ID,
			Firstname: worker.Firstname,
			Lastname:  worker.Lastname,
			Position:  worker.Position,
		}
		Resp = append(Resp, respWorker)
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, requestLang)
	clog.Info(handleName + " success")
}

/////POSITION

func (s *Server) HandlePositionList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandlePositionList"

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

	_, err = s.UserRequirments(ctx, w, r, roles)
	if err != nil {
		eMsg := "UserRequirments error in " + handleName
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.PositionList(ctx)
	if err != nil {
		eMsg := "error in s.c.PositionList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, *data, clog, requestLang)
	clog.Info(handleName + " success")
}
