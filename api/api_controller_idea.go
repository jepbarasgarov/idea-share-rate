package api

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CRITERIA

//////////////////////////////////////////////////////////////////////////////////////////////MONGO////////////////////////////////////////////////////////////////////////

//IDEA

func (api *APIController) IdeaCreate(
	ctx context.Context,
	Idea *models.IdeaCreate,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaCreate",
	})

	worker, err := api.access.Workerget(ctx, Idea.Worker.ID)
	if err != nil {
		eMsg := "error in api.access.Workerget"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if worker == nil {
		eMsg := "worker not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_WORKER)
		return
	}
	Idea.Worker = *worker

	hasRestrctionForWorker, err := api.cache.HasRestrctionForWorker(ctx, worker.Firstname+worker.LastName)
	if err != nil {
		eMsg := "error ocurred while checking if ip has HasRestrctionForWorker"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if hasRestrctionForWorker {
		eMsg := "idea submit rate has not passed yet"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_IP_RESTRICTED)
		return
	}

	genres, err := api.access.GenreList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	genreNum := len(*genres)
	var genreValid bool = false
	for i := 0; i < genreNum; i++ {
		if (*genres)[i] == Idea.Genre {
			genreValid = true
			break
		}
	}

	if !genreValid {
		eMsg := "Genre not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_GENRE)
		return
	}

	mechanicsValid, err := api.access.CheckAllMechanicsArePresent(ctx, Idea.Mechanics)
	if err != nil {
		eMsg := "error in api.access.CheckAllMechanicsArePresent"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if !mechanicsValid {
		eMsg := "Mechanics not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_MECH)
		return
	}

	breakingIndex := -1
	for i, pFile := range Idea.Files {
		path, err := helpers.ProcessFile(ctx, config.Conf.StaticDir, &pFile)
		if err != nil {
			clog.WithError(err).Error("an error ocurred on converter.ProcessFile")
			breakingIndex = i
			break
		}
		sketch := models.SketchSturctInIdea{
			SketchID: primitive.NewObjectID(),
			Path:     path,
			FileName: pFile.FileHeader.Filename,
		}

		Idea.AllFiles = append(Idea.AllFiles, sketch)
	}

	defer func() {
		if err != nil || breakingIndex != -1 {
			for _, sketch := range Idea.AllFiles {
				_ = os.Remove(filepath.Join(config.Conf.StaticDir, sketch.Path))
			}
		}
	}()

	if breakingIndex != -1 {
		eMsg := "error ocurred while writing the files"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	err = api.access.IdeaCreate(ctx, Idea)
	if err != nil {
		eMsg := "error in api.access.IdeaCreate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	_ = api.cache.WorkerRestrictWithExpiry(ctx, worker.Firstname+worker.LastName, time.Minute*5)

	return
}

func (api *APIController) IdeaRate(
	ctx context.Context,
	Rating *models.RateIdeaCritera,
	cu *responses.ActionInfo,
) (item *int, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaRate",
	})

	idea, err := api.access.IdeaGet(ctx, cu, Rating.IdeaID)
	if err != nil {
		eMsg := "error in api.access.GetIdeaByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if idea == nil {
		eMsg := "idea not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_IDEA)
		return
	}

	criterNum := len(idea.CriteriaRates)
	var criteriaExists bool = false
	for i := 0; i < criterNum; i++ {
		if idea.CriteriaRates[i].ID == Rating.Rating.CriteriaID {
			criteriaExists = true
			break
		}
	}

	if !criteriaExists {
		eMsg := "Criteria not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_CRITERIA)
		return
	}

	criter, err := api.access.CriteriaGetByID(ctx, Rating.Rating.CriteriaID)
	if err != nil {
		eMsg := "error in api.access.CriteriaGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if criter == nil {
		eMsg := "Criteria not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_CRITERIA)
		return
	}

	Rating.Rating.CrieteriaName = criter.Name
	Rating.Rating.UserID, _ = primitive.ObjectIDFromHex(cu.ID)

	item, err = api.access.IdeaRate(ctx, Rating)
	if err != nil {
		eMsg := "error in api.access.IdeaRate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) IdeaList(
	ctx context.Context,
	Filter *models.IdeaFilter,
	cu *responses.ActionInfo,
) (item *models.IdeaList, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaList",
	})

	item, err = api.access.IdeaList(ctx, cu, Filter)
	if err != nil {
		eMsg := "error in api.access.IdeaList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

func (api *APIController) IdeaGet(
	ctx context.Context,
	id primitive.ObjectID,
	cu *responses.ActionInfo,
) (item *models.IdeaSpecData, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaGet",
	})

	item, err = api.access.IdeaGet(ctx, cu, id)
	if err != nil {
		eMsg := "error in api.access.IdeaGet"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if item == nil {
		eMsg := "Idea not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_IDEA)
		return
	}
	return
}

func (api *APIController) IdeaDelete(
	ctx context.Context,
	id primitive.ObjectID,
	cu *responses.ActionInfo,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaGet",
	})

	idea, err := api.access.IdeaGet(ctx, cu, id)
	if err != nil {
		eMsg := "error in api.access.GetIdeaByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if idea == nil {
		eMsg := "Idea not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_IDEA)
		return
	}

	err = api.access.IdeaDelete(ctx, id)
	if err != nil {
		eMsg := "error in api.access.IdeaDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) IdeaUpdate(
	ctx context.Context,
	cu *responses.ActionInfo,
	newIdea *models.IdeaUpdate,
) (item *models.IdeaSpecData, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.IdeaUpdate",
	})

	oldIdea, err := api.access.IdeaGet(ctx, cu, newIdea.ID)
	if err != nil {
		eMsg := "error in api.access.GetIdeaByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if oldIdea == nil {
		eMsg := "Idea not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_IDEA)
		return
	}

	worker, err := api.access.Workerget(ctx, newIdea.Worker.ID)
	if err != nil {
		eMsg := "error in api.access.Workerget"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if worker == nil {
		eMsg := "worker not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_WORKER)
		return
	}

	newIdea.Worker = *worker

	genres, err := api.access.GenreList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	genreNum := len(*genres)
	var genreValid bool = false
	for i := 0; i < genreNum; i++ {
		if (*genres)[i] == newIdea.Genre {
			genreValid = true
			break
		}
	}

	if !genreValid {
		eMsg := "Genre not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_GENRE)
		return
	}

	mechanicsValid, err := api.access.CheckAllMechanicsArePresent(ctx, newIdea.Mechanics)
	if err != nil {
		eMsg := "error in api.access.CheckAllMechanicsArePresent"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if !mechanicsValid {
		eMsg := "Mechanics not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_MECH)
		return
	}

	err = api.access.IdeaUpdate(ctx, nil, newIdea)
	if err != nil {
		eMsg := "error in api.access.IdeaUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &models.IdeaSpecData{
		ID:            newIdea.ID,
		Name:          newIdea.Name,
		Genre:         newIdea.Genre,
		Description:   newIdea.Description,
		Date:          newIdea.Date,
		Mechanics:     newIdea.Mechanics,
		Links:         newIdea.Links,
		Worker:        *worker,
		FilePaths:     oldIdea.FilePaths,
		CriteriaRates: oldIdea.CriteriaRates,
		OverallRate:   oldIdea.OverallRate,
	}
	return
}

//GENRE

func (api *APIController) GenreCreate(
	ctx context.Context,
	genre string,
) (item *string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.GenreCreate",
	})
	err = api.access.GenreUpsert(ctx, genre)
	if err != nil {
		eMsg := "error in api.access.GenreUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &genre
	return
}

func (api *APIController) GenreList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.GenreList",
	})
	item, err = api.access.GenreList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

func (api *APIController) GenreDelete(
	ctx context.Context,
	genre string,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.GenreDelete",
	})

	genreList, err := api.access.GenreList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	var genreExist bool = false

	for i := 0; i < len(*genreList); i++ {
		if (*genreList)[i] == genre {
			genreExist = true
		}
	}

	if !genreExist {
		eMsg := "Genre not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_GENRE)
		return
	}

	err = api.access.GenreDelete(ctx, genre)
	if err != nil {
		eMsg := "error in api.access.GenreDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) GenreUpdate(
	ctx context.Context,
	GenreUpdate models.GenreUpdate,
) (item *string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.GenreUpdate",
	})

	genreList, err := api.access.GenreList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	var oldExist, newExist bool = false, false

	for i := 0; i < len(*genreList); i++ {
		if (*genreList)[i] == GenreUpdate.OldGenre {
			oldExist = true
		}
		if (*genreList)[i] == GenreUpdate.NewGenre {
			newExist = true
		}

	}

	if !oldExist {
		eMsg := "Genre not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_GENRE)
		return
	}

	if GenreUpdate.OldGenre == GenreUpdate.NewGenre {
		item = &GenreUpdate.NewGenre
		return
	}

	if newExist {
		eMsg := "Genre is already in use"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_GENRE)
		return
	}

	err = api.access.GenreUpdate(ctx, GenreUpdate)
	if err != nil {
		eMsg := "error in api.access.GenreUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &GenreUpdate.NewGenre
	return
}

//MECHANICS

func (api *APIController) MechanicCreate(
	ctx context.Context,
	mechanic string,
) (item *string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.MechanicCreate",
	})
	err = api.access.MechanicUpsert(ctx, mechanic)
	if err != nil {
		eMsg := "error in api.access.GenreUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &mechanic
	return
}

func (api *APIController) MechanicList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.MechanicList",
	})
	item, err = api.access.MechanicList(ctx)
	if err != nil {
		eMsg := "error in api.access.MechanicList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

func (api *APIController) MechanicDelete(
	ctx context.Context,
	mechanic string,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.MechanicDelete",
	})

	mechList, err := api.access.MechanicList(ctx)
	if err != nil {
		eMsg := "error in api.access.MechanicList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	var mechanicExist bool = false

	for i := 0; i < len(*mechList); i++ {
		if (*mechList)[i] == mechanic {
			mechanicExist = true
		}
	}

	if !mechanicExist {
		eMsg := "mechanic not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_MECH)
		return
	}

	err = api.access.MechanicDelete(ctx, mechanic)
	if err != nil {
		eMsg := "error in api.access.MechanicDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) MechanicUpdate(
	ctx context.Context,
	MechUpdate models.MechanicUpdate,
) (item *string, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.MechanicUpdate",
	})

	MechList, err := api.access.MechanicList(ctx)
	if err != nil {
		eMsg := "error in api.access.GenreUpsert"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	var oldExist, newExist bool = false, false

	for i := 0; i < len(*MechList); i++ {
		if (*MechList)[i] == MechUpdate.OldMech {
			oldExist = true
		}
		if (*MechList)[i] == MechUpdate.NewMech {
			newExist = true
		}

	}

	if !oldExist {
		eMsg := "Mechanics not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_MECH)
		return
	}

	if MechUpdate.OldMech == MechUpdate.NewMech {
		item = &MechUpdate.NewMech
		return
	}

	if newExist {
		eMsg := "Mechanic is already in use"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_MECH)
		return
	}

	err = api.access.MechanicUpdate(ctx, MechUpdate)
	if err != nil {
		eMsg := "error in api.access.MechanicUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &MechUpdate.NewMech
	return
}

//CRITERIA

func (api *APIController) CriteriaCreate(
	ctx context.Context,
	cu *responses.ActionInfo,
	CriteriaName string,
) (item *models.CriteriaSpecData, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.CriteriaCreate",
		"username": cu.Username,
	})

	criteriaID, err := api.access.CriteriaGetByName(ctx, CriteriaName)
	if err != nil {
		eMsg := "error in api.access.CriteriaGetByName"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if criteriaID != nil {
		eMsg := "Criteria name is in use"
		clog.Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_CRITERIA)
		return
	}

	item, err = api.access.CriteriaCreate(ctx, CriteriaName)
	if err != nil {
		eMsg := "error in api.access.CriteriaCreate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) CriteriaUpdate(
	ctx context.Context,
	cu *responses.ActionInfo,
	criteria *models.CriteriaUpdate,
) (err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.CriteriaUpdate",
		"username": cu.Username,
	})

	criterGotByID, err := api.access.CriteriaGetByID(ctx, criteria.ID)
	if err != nil {
		eMsg := "error in api.access.CriteriaGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if criterGotByID == nil {
		eMsg := "Criteria not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_CRITERIA)
		return
	}

	criterIDGotByName, err := api.access.CriteriaGetByName(ctx, criteria.Name)
	if err != nil {
		eMsg := "error in api.access.CriteriaGetByName"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if criterIDGotByName != nil && *criterIDGotByName != criterGotByID.ID {
		eMsg := "Criteria name is in use"
		clog.Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_CRITERIA)
		return
	}

	err = api.access.CriteriaUpdate(ctx, criteria)
	if err != nil {
		eMsg := "error in api.access.CriteriaUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) CriteriaList(
	ctx context.Context,
) (item *[]models.CriteriaSpecData, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.CriteriaList",
	})
	item, err = api.access.CriteriaList(ctx)
	if err != nil {
		eMsg := "error in api.access.CriteriaList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

func (api *APIController) CriteriaDelete(
	ctx context.Context,
	ID primitive.ObjectID,
	cu *responses.ActionInfo,
) (err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.CriteriaDelete",
		"username": cu.Username,
	})

	cirteriaGotByID, err := api.access.CriteriaGetByID(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.CriteriaGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if cirteriaGotByID == nil {
		eMsg := "Criteria not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_CRITERIA)
		return
	}

	numOfRates, err := api.access.CountCriteriaRates(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.CountCriteriaRates"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if numOfRates > 0 {
		eMsg := "Criteria has some rates"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_CRITERIA_HAS_RATES)
		return
	}

	err = api.access.CriteriaDelete(ctx, nil, ID)
	if err != nil {
		eMsg := "error in api.access.CriteriaDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}
