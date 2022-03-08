package web

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//IDEA

func (s *Server) HandleIdeaList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaList"

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

	var Filter models.IdeaFilter

	Filter.UserID, _ = primitive.ObjectIDFromHex(cu.ID)

	if len(r.FormValue("name")) != 0 {
		name := r.FormValue("name")
		Filter.Name = &name
	}

	woID, err := primitive.ObjectIDFromHex(r.FormValue("worker_id"))
	if err == nil {
		Filter.WorkerID = &woID
	}

	dateStart, err := helpers.ChangeStringToDate(r.FormValue("start"))
	if err == nil && !dateStart.IsZero() {
		d := r.FormValue("start")
		Filter.BeginDate = &d
	}

	dateEnd, err := helpers.ChangeStringToDate(r.FormValue("end"))
	if err == nil && !dateEnd.IsZero() {
		d := r.FormValue("end")
		Filter.BeginDate = &d
	}

	if len(r.FormValue("genre")) != 0 {
		genre := r.FormValue("genre")
		Filter.Genre = &genre
	}

	mechanics := make([]string, 0)
	err = json.Unmarshal([]byte(r.FormValue("mechanics")), &mechanics)
	if err == nil {
		Filter.Mechanics = &mechanics
	}

	condition, err := helpers.ConvertStringToIdeaCondition(r.FormValue("condition"))
	if err == nil {
		Filter.Condition = &condition
	}

	Filter.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		Filter.Limit = 50
	}

	Filter.Offset, err = strconv.Atoi(r.FormValue("skip"))
	if err != nil {
		Filter.Offset = 0
	}

	data, err := s.c.IdeaList(ctx, &Filter, cu)
	if err != nil {
		eMsg := "error in s.c.IdeaList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleIdeaGet(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaGet"

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

	x := mux.Vars(r)["id"]
	id, err := uuid.FromString(x)
	if err != nil {
		emsg := "IdeaID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.IdeaGet(ctx, id.String(), cu)
	if err != nil {
		eMsg := "error in s.c.IdeaGet"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Resp := responses.IdeaSpecData{
		ID:          data.ID,
		Name:        data.Name,
		Genre:       data.Genre,
		Description: data.Description,
		Worker: responses.WorkerLightData{
			ID:        data.Worker.ID,
			Firstname: data.Worker.Firstname,
			Lastname:  data.Worker.Lastname,
			Position:  data.Worker.Position,
		},
		Date:          data.Date,
		Mechanics:     data.Mechanics,
		Links:         data.Links,
		FilePaths:     data.FilePaths,
		CriteriaRates: data.CriteriaRates,
		OverallRate:   data.OverallRate,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleIdeaListGetPdf(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaListGetPdf"

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

	var Filter models.IdeaFilter
	Filter.UserID, _ = primitive.ObjectIDFromHex(cu.ID)

	if len(r.FormValue("name")) != 0 {
		name := r.FormValue("name")
		Filter.Name = &name
	}

	woID, err := primitive.ObjectIDFromHex(r.FormValue("worker_id"))
	if err == nil {
		Filter.WorkerID = &woID
	}

	dateStart, err := helpers.ChangeStringToDate(r.FormValue("start"))
	if err == nil && !dateStart.IsZero() {
		d := r.FormValue("start")
		Filter.BeginDate = &d
	}

	dateEnd, err := helpers.ChangeStringToDate(r.FormValue("end"))
	if err == nil && !dateEnd.IsZero() {
		d := r.FormValue("end")
		Filter.BeginDate = &d
	}

	if len(r.FormValue("genre")) != 0 {
		genre := r.FormValue("genre")
		Filter.Genre = &genre
	}

	mechanics := make([]string, 0)
	err = json.Unmarshal([]byte(r.FormValue("mechanics")), &mechanics)
	if err == nil {
		Filter.Mechanics = &mechanics
	}

	condition, err := helpers.ConvertStringToIdeaCondition(r.FormValue("condition"))
	if err == nil {
		Filter.Condition = &condition
	}

	Filter.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		Filter.Limit = 50
	}

	Filter.Offset, err = strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		Filter.Offset = 0
	}

	data, err := s.c.IdeaList(ctx, &Filter, cu)
	if err != nil {
		eMsg := "error in s.c.IdeaList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	// Resp := responses.IdeaList{}
	// Resp.Total = data.Total
	// Resp.LastSubmitted = data.LastSubmitted
	// Resp.Result = make([]responses.IdeaLightData, 0)

	// for _, idea := range data.Result {
	// 	respIdea := responses.IdeaLightData{
	// 		ID:   idea.ID,
	// 		Name: idea.Name,
	// 		Worker: responses.WorkerLightData{
	// 			ID:        idea.Worker.ID,
	// 			Firstname: idea.Worker.Firstname,
	// 			Lastname:  idea.Worker.Lastname,
	// 			Position:  idea.Worker.Position,
	// 		},
	// 		Date:        idea.Date,
	// 		Description: idea.Description,
	// 		IsItNew:     idea.IsItNew,
	// 		FilePath:    idea.FilePath,
	// 		OverallRate: idea.OverallRate,
	// 	}
	// 	Resp.Result = append(Resp.Result, respIdea)
	// }

	renderInfo, err := json.Marshal(data.Result)
	if err != nil {
		eMsg := "error in marshalling render info"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	response, err := http.Post(fmt.Sprintf("%s/render-list", config.Conf.Renderer), "application/json", bytes.NewBuffer(renderInfo))
	if err != nil {
		eMsg := "error in renderer post"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		eMsg := "error in renderer post"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func (s *Server) HandleIdeaGetPdf(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaGetPdf"

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

	x := mux.Vars(r)["id"]
	id, err := uuid.FromString(x)
	if err != nil {
		emsg := "IdeaID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.IdeaGet(ctx, id.String(), cu)
	if err != nil {
		eMsg := "error in s.c.IdeaGet"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Resp := responses.IdeaSpecData{
		ID:          data.ID,
		Name:        data.Name,
		Genre:       data.Genre,
		Description: data.Description,
		Worker: responses.WorkerLightData{
			ID:        data.Worker.ID,
			Firstname: data.Worker.Firstname,
			Lastname:  data.Worker.Lastname,
			Position:  data.Worker.Position,
		},
		Date:          data.Date,
		Mechanics:     data.Mechanics,
		Links:         data.Links,
		FilePaths:     data.FilePaths,
		CriteriaRates: data.CriteriaRates,
		OverallRate:   data.OverallRate,
	}

	renderInfo, err := json.Marshal(Resp)
	if err != nil {
		eMsg := "error in marshalling render info"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	response, err := http.Post(fmt.Sprintf("%s/render-item", config.Conf.Renderer), "application/json", bytes.NewBuffer(renderInfo))
	if err != nil {
		eMsg := "error in renderer post"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		eMsg := "error in renderer post"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func (s *Server) HandleIdeaDelete(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaDelete"

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

	x := mux.Vars(r)["id"]
	id, err := uuid.FromString(x)
	if err != nil {
		emsg := "IdeaID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = s.c.IdeaDelete(ctx, id.String(), cu)
	if err != nil {
		eMsg := "error in s.c.IdeaDelete"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleIdeaUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaUpdate"

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

	var IdeaUpdate models.IdeaUpdate

	_, err = uuid.FromString(mux.Vars(r)["id"])
	if err != nil {
		emsg := "IdeaID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	IdeaUpdate.ID = mux.Vars(r)["id"]

	if len(r.FormValue("name")) == 0 || len(r.FormValue("name")) > 256 {
		emsg := "Idea name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	IdeaUpdate.Name = r.FormValue("name")

	_, err = uuid.FromString(r.FormValue("worker_id"))
	if err != nil {
		emsg := "Worker ID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	IdeaUpdate.WorkerID = r.FormValue("worker_id")

	IdeaUpdate.Date, err = helpers.ChangeStringToDate(r.FormValue("date"))
	if err != nil || IdeaUpdate.Date.IsZero() {
		eMsg := "There is no such date"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	if len(r.FormValue("genre")) == 0 || len(r.FormValue("genre")) > 256 {
		emsg := "Idea genre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	IdeaUpdate.Genre = r.FormValue("genre")

	mechanics := make([]string, 0)
	err = json.Unmarshal([]byte(r.FormValue("mechanics")), &mechanics)
	if err != nil {
		eMsg := "Mechanics must be string"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	sort.Strings(mechanics)
	mechs := len(mechanics)
	if mechs > 0 {
		if len(mechanics[0]) < 256 && len(mechanics[0]) > 0 {
			IdeaUpdate.Mechanics = append(IdeaUpdate.Mechanics, mechanics[0])
		}
		if mechs > 1 {
			for i := 1; i < mechs; i++ {
				if mechanics[i] != mechanics[i-1] {
					if len(mechanics[i]) < 256 && len(mechanics[i]) > 0 {
						IdeaUpdate.Mechanics = append(IdeaUpdate.Mechanics, mechanics[i])
					}
				}
			}
		}
	}

	if len(IdeaUpdate.Mechanics) == 0 {
		eMsg := "There must be related mechanics"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	IdeaUpdate.Description = r.FormValue("description")

	err = json.Unmarshal([]byte(r.FormValue("links")), &IdeaUpdate.Links)
	if err != nil {
		eMsg := "Links must contain label and url"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	// data, err := s.c.IdeaUpdate(ctx, cu, &IdeaUpdate)
	// if err != nil {
	// 	eMsg := "error in s.c.IdeaUpdate"
	// 	clog.WithError(err).Error(eMsg)
	// 	errs.SendResponse(w, err, nil, clog, requestLang)
	// 	return
	// }

	// Resp := responses.IdeaSpecData{
	// 	ID:          data.ID,
	// 	Name:        data.Name,
	// 	Genre:       data.Genre,
	// 	Description: data.Description,
	// 	Worker: responses.WorkerLightData{
	// 		ID:        data.Worker.ID,
	// 		Firstname: data.Worker.Firstname,
	// 		Lastname:  data.Worker.Lastname,
	// 		Position:  data.Worker.Position,
	// 	},
	// 	Date:          data.Date,
	// 	Mechanics:     data.Mechanics,
	// 	Links:         data.Links,
	// 	FilePaths:     data.FilePaths,
	// 	CriteriaRates: data.CriteriaRates,
	// 	OverallRate:   data.OverallRate,
	// }

	// err = responses.ErrOK
	// errs.SendResponse(w, err, Resp, clog, requestLang)
	// clog.Info(handleName + " success")
}

//GENRE

func (s *Server) HandleGenreUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleGenreUpdate"

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

	var GenreUpdate models.GenreUpdate
	GenreUpdate.OldGenre = r.FormValue("old_genre")
	if len(GenreUpdate.OldGenre) == 0 || len(GenreUpdate.OldGenre) > 256 {
		emsg := "OldGenre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	GenreUpdate.NewGenre = r.FormValue("new_genre")
	if len(GenreUpdate.NewGenre) == 0 || len(GenreUpdate.NewGenre) > 256 {
		emsg := "NewGenre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.GenreUpdate(ctx, GenreUpdate)
	if err != nil {
		eMsg := "error in s.c.GenreUpdate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, requestLang)
	clog.Info(handleName + " success")
}

//MECHANICS

func (s *Server) HandleMechanicUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleMechanicUpdate"

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

	var MechUpdate models.MechanicUpdate
	MechUpdate.OldMech = r.FormValue("old_mechanic")
	if len(MechUpdate.OldMech) == 0 || len(MechUpdate.OldMech) > 256 {
		emsg := "OldMechanic name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	MechUpdate.NewMech = r.FormValue("new_mechanic")
	if len(MechUpdate.NewMech) == 0 || len(MechUpdate.NewMech) > 256 {
		emsg := "NewMechanic name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.MechanicUpdate(ctx, MechUpdate)
	if err != nil {
		eMsg := "error in s.c.MechanicUpdate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, requestLang)
	clog.Info(handleName + " success")
}

//CRITERIA

// func (s *Server) HandleCriteriaDelete(w http.ResponseWriter, r *http.Request) {
// 	handleName := "HandleCriteriaDelete"

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

// 	ID, err := uuid.FromString(mux.Vars(r)["id"])
// 	if err != nil {
// 		eMsg := "criteria ID is not compatible"
// 		clog.Error(eMsg)
// 		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
// 		errs.SendResponse(w, err, nil, clog, cu.Language)
// 		return
// 	}

// 	err = s.c.CriteriaDelete(ctx, ID.String(), cu)
// 	if err != nil {
// 		eMsg := "error in s.c.CriteriaDelete"
// 		clog.WithError(err).Error(eMsg)
// 		errs.SendResponse(w, err, nil, clog, cu.Language)
// 		return
// 	}

// 	err = responses.ErrOK
// 	errs.SendResponse(w, err, nil, clog, cu.Language)

// 	clog.Info(handleName + " success")
// }

////////////////////////////////////////////////////////////////////////////////MONGO////////////////////////////////////////////////////////////////////////////////////

//IDEA

func (s *Server) HandleIdeaCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleIdeaCreate"

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

	var Idea models.IdeaCreate
	if len(r.FormValue("name")) == 0 || len(r.FormValue("name")) > 256 {
		emsg := "Idea name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	Idea.Name = r.FormValue("name")

	Idea.Worker.ID, err = primitive.ObjectIDFromHex(r.FormValue("worker_id"))
	if err != nil {
		emsg := "Worker ID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Idea.Date, err = helpers.ChangeStringToDate(r.FormValue("date"))
	if err != nil || Idea.Date.IsZero() {
		eMsg := "There is no such date"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	if len(r.FormValue("genre")) == 0 || len(r.FormValue("genre")) > 256 {
		emsg := "Idea genre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	Idea.Genre = r.FormValue("genre")

	mechanics := make([]string, 0)
	err = json.Unmarshal([]byte(r.FormValue("mechanics")), &mechanics)
	if err != nil {
		eMsg := "Mechanics must be string"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	sort.Strings(mechanics)
	mechs := len(mechanics)
	if mechs > 0 {
		if len(mechanics[0]) < 256 && len(mechanics[0]) > 0 {
			Idea.Mechanics = append(Idea.Mechanics, mechanics[0])
		}
		if mechs > 1 {
			for i := 1; i < mechs; i++ {
				if mechanics[i] != mechanics[i-1] {
					if len(mechanics[i]) < 256 && len(mechanics[i]) > 0 {
						Idea.Mechanics = append(Idea.Mechanics, mechanics[i])
					}
				}
			}
		}
	}

	if len(Idea.Mechanics) == 0 {
		eMsg := "There must be related mechanics"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Idea.Description = r.FormValue("description")

	err = json.Unmarshal([]byte(r.FormValue("links")), &Idea.Links)
	if err != nil {
		eMsg := "Links must contain label and url"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	_, _, err = r.FormFile("files")
	if err == nil {
		fhs := r.MultipartForm.File["files"]
		for _, fh := range fhs {
			f, err := fh.Open()
			if err != nil {
				eMsg := "an error occurred on fh.Open"
				clog.WithError(err).Error(eMsg)
				err = errs.NewHttpErrorInternalError(errs.ERR_IE)
				errs.SendResponse(w, err, nil, clog, requestLang)
				return
			}

			Idea.Files = append(Idea.Files, models.ParsedFile{
				File:        f,
				FileHeader:  fh,
				ContentType: helpers.GetContentType(fh.Filename),
			})
		}

		for _, pFile := range Idea.Files {
			if !helpers.CheckFileContentType(pFile.ContentType, config.Conf.IdeaFile.ContentTypes) {
				eMsg := fmt.Sprintf("content type isn't allowed: %s", pFile.FileHeader.Filename)
				clog.WithError(err).Error(eMsg)
				err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
				errs.SendResponse(w, err, nil, clog, requestLang)
				return
			}

			if pFile.FileHeader.Size > config.Conf.IdeaFile.MaxSize {
				eMsg := fmt.Sprintf("File size is too large file_name = %s", pFile.FileHeader.Filename)
				clog.WithError(err).Error(eMsg)
				err = errs.NewHttpErrorFileTooLarge(errs.ERR_FL)
				errs.SendResponse(w, err, nil, clog, requestLang)
				return
			}

			if len(pFile.FileHeader.Filename) > 256 {
				eMsg := fmt.Sprintf("Filename is too big: %s", pFile.FileHeader.Filename)
				clog.WithError(err).Error(eMsg)
				err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
				errs.SendResponse(w, err, nil, clog, requestLang)
				return
			}
		}
	}

	err = s.c.IdeaCreate(ctx, &Idea)
	if err != nil {
		eMsg := "error in s.c.IdeaCreate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleRateIdea(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleRateIdea"

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

	var Rating models.RateIdeaCritera

	Rating.IdeaID, err = primitive.ObjectIDFromHex(r.FormValue("idea_id"))
	if err != nil {
		emsg := "Idea ID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}
	Rating.Rating.CriteriaID, err = primitive.ObjectIDFromHex(r.FormValue("criteria_id"))
	if err != nil {
		emsg := "Criteria ID is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Rating.Rating.Rate, err = strconv.Atoi(r.FormValue("rate"))
	if err != nil || Rating.Rating.Rate <= 0 || Rating.Rating.Rate >= 11 {
		emsg := "Rate value is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.IdeaRate(ctx, &Rating, cu)
	if err != nil {
		eMsg := "error in s.c.IdeaCreate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Resp := responses.OverAllRate{
		Rate: *data,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, requestLang)
	clog.Info(handleName + " success")
}

//Genre
func (s *Server) HandleGenreCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleGenreCreate"

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
	genre := r.FormValue("genre")
	if len(genre) == 0 || len(genre) > 256 {
		emsg := "Genre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.GenreCreate(ctx, genre)
	if err != nil {
		eMsg := "error in s.c.GenreCreate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleGenreList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleGenreList"

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

	data, err := s.c.GenreList(ctx)
	if err != nil {
		eMsg := "error in s.c.GenreList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, *data, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleGenreDelete(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleGenreDelete"

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
	genre := r.FormValue("genre")
	if len(genre) == 0 || len(genre) > 256 {
		emsg := "Genre name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = s.c.GenreDelete(ctx, genre)
	if err != nil {
		eMsg := "error in s.c.GenreDelete"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, requestLang)
	clog.Info(handleName + " success")
}

//Mechanics

func (s *Server) HandleMechanicsCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleMechanicsCreate"

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
	mechanic := r.FormValue("mechanic")
	if len(mechanic) == 0 || len(mechanic) > 256 {
		emsg := "mechanic name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	data, err := s.c.MechanicCreate(ctx, mechanic)
	if err != nil {
		eMsg := "error in s.c.MechanicCreate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, data, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleMechanicList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleMechanicList"

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

	data, err := s.c.MechanicList(ctx)
	if err != nil {
		eMsg := "error in s.c.MechanicList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, *data, clog, requestLang)
	clog.Info(handleName + " success")
}

func (s *Server) HandleMechanicDelete(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleMechanicDelete"

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
	mechanic := r.FormValue("mechanic")
	if len(mechanic) == 0 || len(mechanic) > 256 {
		emsg := "mechanic name is not compatible"
		clog.WithError(err).Error(emsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = s.c.MechanicDelete(ctx, mechanic)
	if err != nil {
		eMsg := "error in s.c.MechanicDelete"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, nil, clog, requestLang)
	clog.Info(handleName + " success")
}

//Criteria

func (s *Server) HandleCriteriaCreate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleCriteriaCreate"

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

	CriteriaName := r.FormValue("name")
	if len(CriteriaName) == 0 || len(CriteriaName) > 256 {
		eMsg := "Criteria name length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	data, err := s.c.CriteriaCreate(ctx, cu, CriteriaName)
	if err != nil {
		eMsg := "error in s.c.CriteriaCreate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	Resp := responses.CriteriaSpecData{
		ID:   data.ID,
		Name: data.Name,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleCriteriaUpdate(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleCriteriaUpdate"

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
	var criteria models.CriteriaUpdate

	criteria.ID, err = primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		eMsg := "Criteria ID is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	criteria.Name = r.FormValue("name")
	if len(criteria.Name) == 0 || len(criteria.Name) > 256 {
		eMsg := "Criteria name length is not compatible"
		clog.Error(eMsg)
		err = errs.NewHttpErrorBadRequest(errs.ERR_BR)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	err = s.c.CriteriaUpdate(ctx, cu, &criteria)
	if err != nil {
		eMsg := "error in s.c.CriteriaUpdate"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, cu.Language)
		return
	}

	Resp := responses.CriteriaSpecData{
		ID:   criteria.ID.Hex(),
		Name: criteria.Name,
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, cu.Language)

	clog.Info(handleName + " success")
}

func (s *Server) HandleCriteriaList(w http.ResponseWriter, r *http.Request) {
	handleName := "HandleCriteriaList"

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

	data, err := s.c.CriteriaList(ctx)
	if err != nil {
		eMsg := "error in s.c.CriteriaList"
		clog.WithError(err).Error(eMsg)
		errs.SendResponse(w, err, nil, clog, requestLang)
		return
	}

	Resp := make([]responses.CriteriaSpecData, 0)
	for _, criter := range *data {
		respCriter := responses.CriteriaSpecData{
			ID:   criter.ID,
			Name: criter.Name,
		}

		Resp = append(Resp, respCriter)
	}

	err = responses.ErrOK
	errs.SendResponse(w, err, Resp, clog, requestLang)
	clog.Info(handleName + " success")
}
