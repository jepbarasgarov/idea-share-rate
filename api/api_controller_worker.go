package api

import (
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/////////////////////////////////////////////////MONGO///////////////////////////////////////////////////////////////////////

func (api *APIController) WorkerCreate(
	ctx context.Context,
	worker *models.WorkerCreate,
	cu *responses.ActionInfo,
) (item *models.WorkerBsonModelInIdea, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.WorkerCreate",
		"username": cu.Username,
	})

	err = api.access.PositionUpsert(ctx, worker.Position)
	if err != nil {
		eMsg := "error in api.access.PositionUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item, err = api.access.WorkerCreate(ctx, worker)
	if err != nil {
		eMsg := "error in api.access.WorkerCreate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) WorkerGetByID(
	ctx context.Context,
	ID primitive.ObjectID,
	cu *responses.ActionInfo,
) (item *models.WorkerBsonModelInIdea, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.WorkerGetByID",
	})

	item, err = api.access.Workerget(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.Workerget"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if item == nil {
		eMsg := "Worker not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_WORKER)
		return
	}
	return
}

func (api *APIController) WorkerDelete(
	ctx context.Context,
	ID primitive.ObjectID,
	cu *responses.ActionInfo,
) (err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.WorkerDelete",
		"username": cu.Username,
	})

	wrkr, err := api.access.Workerget(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.Workerget"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if wrkr == nil {
		eMsg := "worker not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_WORKER)
		return
	}
	numOfIdeas, err := api.access.CountWorkersIdea(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.CountWorkersIdea"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if numOfIdeas > 0 {
		eMsg := "worker has some ideas"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_WORKER_HAS_IDEA)
		return
	}

	err = api.access.WorkerDelete(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.WorkerDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) WorkerUpdate(
	ctx context.Context,
	worker *models.WorkerUpdate,
	cu *responses.ActionInfo,
) (item *models.WorkerBsonModelInIdea, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.WorkerUpdate",
		"username": cu.Username,
	})

	wrkr, err := api.access.Workerget(ctx, worker.ID)
	if err != nil {
		eMsg := "error in api.access.Workerget"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if wrkr == nil {
		eMsg := "worker not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_WORKER)
		return
	}

	err = api.access.PositionUpsert(ctx, worker.Position)
	if err != nil {
		eMsg := "error in api.access.PositionUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item, err = api.access.WorkerUpdate(ctx, worker)
	if err != nil {
		eMsg := "error in api.access.WorkerUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) WorkerAutoCompleteList(
	ctx context.Context,
) (item *[]models.WorkerBsonModelInIdea, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.WorkerAutoCompleteList",
	})

	item, err = api.access.WorkerAutocompleteList(ctx)
	if err != nil {
		eMsg := "error in api.access.WorkerAutocompleteList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

// POSITION

func (api *APIController) PositionList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.PositionList",
	})
	item, err = api.access.PositionList(ctx)
	if err != nil {
		eMsg := "error in api.access.PositionList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}
