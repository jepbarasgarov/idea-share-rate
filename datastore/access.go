package datastore

import (
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"

	"github.com/jackc/pgx/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Access interface {
	/////////////////////////////////////////////////////////////////-------- USER
	UserGetByUsername(
		ctx context.Context,
		username string,
	) (item *models.UserSpecDataBson, err error)

	UserCreate(
		ctx context.Context,
		user *models.UserCreate,
	) (item *models.UserSpecDataBson, err error)

	UserGetPasswordByID(
		ctx context.Context,
		id primitive.ObjectID,
	) (pwdHash string, err error)

	UserUpdateOwnPassword(
		ctx context.Context,
		cu *responses.ActionInfo,
		newPassword string,
	) (err error)

	UserGetByID(
		ctx context.Context,
		id primitive.ObjectID,
	) (item *models.UserSpecDataBson, err error)

	AdminUpdatePassword(
		ctx context.Context,
		cu *responses.ActionInfo,
		userid primitive.ObjectID,
		newPassword string,
	) (err error)

	UserUpdate(
		ctx context.Context,
		user *models.UserUpdate,
	) (item *models.UserSpecDataBson, err error)

	UserDelete(
		ctx context.Context,
		id string,
	) (err error)

	UserAutocompleteList(
		ctx context.Context,
	) (item *[]models.UserLightData, err error)

	///////////////////////////////////////////////////////////////-----------WORKER
	WorkerAutocompleteList(
		ctx context.Context,
	) (item *[]models.WorkerBsonModelInIdea, err error)

	Workerget(
		ctx context.Context,
		ID primitive.ObjectID,
	) (item *models.WorkerBsonModelInIdea, err error)

	WorkerCreate(
		ctx context.Context,
		worker *models.WorkerCreate,
	) (item *models.WorkerBsonModelInIdea, err error)

	WorkerUpdate(
		ctx context.Context,
		worker *models.WorkerUpdate,
	) (item *models.WorkerBsonModelInIdea, err error)

	CountWorkersIdea(
		ctx context.Context,
		ID primitive.ObjectID,
	) (item int, err error)

	WorkerDelete(
		ctx context.Context,
		ID primitive.ObjectID,
	) (err error)

	////////////////////////////////////////////////////////////////-----------POSITION
	PositionUpsert(
		ctx context.Context,
		PositionName string,
	) (err error)

	PositionList(
		ctx context.Context,
	) (item *[]string, err error)

	////////////////////////////////////////////////////////////////------------GENRES
	GenreList(
		ctx context.Context,
	) (item *[]string, err error)

	GenreUpsert(
		ctx context.Context,
		GenreName string,
	) (err error)

	GenreUpdate(
		ctx context.Context,
		GenreUpdate models.GenreUpdate,
	) (err error)

	GenreDelete(
		ctx context.Context,
		GenreName string,
	) (err error)

	/////////////////////////////////////////////////////////////--------------MECHANICS
	MechanicList(
		ctx context.Context,
	) (item *[]string, err error)

	MechanicUpsert(
		ctx context.Context,
		Mechanics string,
	) (err error)

	CheckAllMechanicsArePresent(
		ctx context.Context,
		mechList []string,
	) (item bool, err error)

	MechanicUpdate(
		ctx context.Context,
		MechUpdate models.MechanicUpdate,
	) (err error)

	MechanicDelete(
		ctx context.Context,
		MechName string,
	) (err error)

	////////////////////////////////////////////////////////////----------------IDEA
	IdeaCreate(
		ctx context.Context,
		Idea *models.IdeaCreate,
	) (err error)

	IdeaList(

		ctx context.Context,
		cu *responses.ActionInfo,
		Filter *models.IdeaFilter,
	) (item *models.IdeaList, err error)

	IdeaGet(
		ctx context.Context,
		cu *responses.ActionInfo,
		ID primitive.ObjectID,
	) (item *models.IdeaSpecData, err error)

	IdeaDelete(
		ctx context.Context,
		ID primitive.ObjectID,
	) (err error)

	IdeaUpdate(
		ctx context.Context,
		pTx pgx.Tx,
		NewIdea *models.IdeaUpdate,
	) (err error)

	IdeaRate(
		ctx context.Context,
		Rating *models.RateIdeaCritera,
	) (item *int, err error)

	/////////////////////////////////////////////////////////////---------------CRITERIA
	CriteriaGetByName(
		ctx context.Context,
		criteriaName string,
	) (item *primitive.ObjectID, err error)

	CriteriaCreate(
		ctx context.Context,
		CriteriaName string,
	) (item *models.CriteriaSpecData, err error)

	CriteriaGetByID(
		ctx context.Context,
		ID primitive.ObjectID,
	) (item *models.CriteriaSpecData, err error)

	CriteriaUpdate(
		ctx context.Context,
		criter *models.CriteriaUpdate,
	) (err error)

	CountCriteriaRates(
		ctx context.Context,
		ID primitive.ObjectID,
	) (item int, err error)

	CriteriaDelete(
		ctx context.Context,
		pTx pgx.Tx,
		ID primitive.ObjectID,
	) (err error)

	CriteriaList(
		ctx context.Context,
	) (item *[]models.CriteriaSpecData, err error)
}
