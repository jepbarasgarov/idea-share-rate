package datastore

import (
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sqlCreateIdea       = `INSERT INTO tbl_idea(name, worker_id, date, genre, mechanics, description) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	sqlUpdateIdea       = `UPDATE tbl_idea SET name = $1, worker_id = $2, date = $3, genre = $4, mechanics =$5 , description = $6 WHERE id = $7`
	sqlDeleteIdea       = `DELETE FROM tbl_idea WHERE id = $1`
	sqlCreateLinkIdea   = `INSERT INTO tbl_link(label, link, idea_id) VALUES($1, $2, $3)`
	sqlDeleteIdeaLinks  = `DELETE FROM tbl_link WHERE idea_id = $1`
	sqlCreateSketchPath = `INSERT INTO tbl_sketch(name, idea_id, file_path, place) VALUES($1, $2, $3, $4)`
	sqlGetIdeaList      = `SELECT idea.id, idea.name, idea.date, idea.description, idea.worker_id,  worker.firstname, worker.lastname, worker.position, sketch.file_path, userrel.mark FROM tbl_idea idea
						   INNER JOIN tbl_worker worker ON idea.worker_id = worker.id
						   LEFT JOIN tbl_sketch sketch ON idea.id = sketch.idea_id 
						   LEFT JOIN tbl_user_idea_rel userrel ON userrel.user_id = $1 AND userrel.idea_id = idea.id
						   WHERE (sketch.place = 1 OR sketch.place IS NULL)`
	sqlSelectLastIdeaSubmittedDateWorker = `SELECT create_ts FROM tbl_idea WHERE worker_id = $1 ORDER BY create_ts DESC LIMIT 1 OFFSET 0`
	sqlCountIdea                         = `SELECT COUNT(*) FROM tbl_idea WHERE 1 = $1 `
	sqlGetOverAllRateIdea                = `SELECT COUNT(*) , COALESCE(SUM(rate), 0) FROM tbl_idea_rate WHERE idea_id = $1`
	sqlGetIdeaByID                       = `SELECT idea.id, idea.name, idea.date, idea.description, idea.genre, idea.mechanics, idea.worker_id,  worker.firstname, worker.lastname, worker.position FROM tbl_idea idea 
 	 						 INNER JOIN tbl_worker worker ON idea.worker_id = worker.id WHERE idea.id = $1`
	sqlGetIdeaLinks         = `SELECT label, link FROM tbl_link WHERE idea_id = $1`
	sqlGetIdeaSketchPaths   = `SELECT id, name, file_path FROM tbl_sketch WHERE idea_id = $1`
	sqlGetRatesOfUserToIdea = `SELECT crit.id, crit.name ,rate.rate FROM tbl_criteria crit 
							  LEFT JOIN tbl_idea_rate rate ON crit.id = rate.criteria_id AND rate.user_id = $1 AND rate.idea_id = $2`
	sqlRateIdea          = `INSERT INTO tbl_idea_rate(idea_id, criteria_id, user_id, rate) VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT user_rate DO UPDATE SET rate = $4, update_ts = $5`
	sqlUpsertIdeaUserRel = `INSERT INTO tbl_user_idea_rel(user_id, idea_id) VALUES($1, $2) ON CONFLICT ON CONSTRAINT user_idea_rel DO NOTHING`

	sqlSelectGenresList        = `SELECT ARRAY(SELECT name FROM tbl_genre)`
	sqlUpsertGenre             = `INSERT INTO tbl_genre(name) VALUES($1) ON CONFLICT ON CONSTRAINT genre_unique DO NOTHING`
	sqlUpdateGenreName         = `UPDATE tbl_genre SET name = $1 WHERE name = $2`
	sqlUpdateAllGenreNamesIdea = `UPDATE tbl_idea SET genre = $1 WHERE genre = $2`
	sqlDeleteGenre             = `DELETE FROM tbl_genre WHERE name = $1`

	sqlSelectMechanicssList       = `SELECT ARRAY(SELECT name FROM tbl_mechanic)`
	sqlUpsertMechanic             = `INSERT INTO tbl_mechanic(name) VALUES ($1) ON CONFLICT ON CONSTRAINT mechanic_unique DO NOTHING`
	sqlUpdateMechanic             = `UPDATE tbl_mechanic SET name = $1 WHERE name = $2`
	sqlUpdateAllmechanicNamesIdea = `UPDATE tbl_idea SET mechanics = array_replace(mechanics , $1, $2)`
	sqlCheckMechanicsArePresent   = `SELECT $1 <@ ARRAY(SELECT name FROM tbl_mechanic)`
	sqlDeleteMechanic             = `DELETE FROM tbl_mechanic WHERE name = $1`

	sqlgetCriteriaByName  = `SELECT id FROM tbl_criteria WHERE name = $1`
	sqlgetCriteriaByID    = `SELECT id, name FROM tbl_criteria WHERE id = $1`
	sqlCreateCriteria     = `INSERT INTO tbl_criteria(name) VALUES($1) RETURNING id`
	sqlUpdateCriteria     = `UPDATE tbl_criteria SET name = $1, update_ts = $2 WHERE id = $3`
	sqlCountCriteriaRates = `SELECT COUNT(*) FROM tbl_idea_rate WHERE criteria_id = $1`
	sqlDeleteCriteria     = `DELETE FROM tbl_criteria WHERE id = $1`
	sqlSelectcriteriaList = `SELECT id, name FROM tbl_criteria`
)

//IDEA

func (d *PgAccess) IdeaList(
	ctx context.Context,
	cu *responses.ActionInfo,
	Filter *models.IdeaFilter,
) (item *models.IdeaList, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaList",
	})

	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {
		defer func() {
			if err != nil {
				item = nil
			}
		}()

		item = &models.IdeaList{}
		item.Result = make([]models.IdeaLightData, 0)

		sqlGet := sqlGetIdeaList
		sqlCount := sqlCountIdea

		ParamsForGet := make([]interface{}, 0)
		ParamsForGet = append(ParamsForGet, cu.ID)
		ParamsForCount := make([]interface{}, 0)
		ParamsForCount = append(ParamsForCount, 1)

		init := 2
		if Filter.WorkerID != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.worker_id = $` + strOrder
			sqlCount += ` AND worker_id = $` + strOrder
			ParamsForGet = append(ParamsForGet, *Filter.WorkerID)
			ParamsForCount = append(ParamsForCount, *Filter.WorkerID)
			init++
		}
		if Filter.Name != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.name ILIKE $` + strOrder
			sqlCount += ` AND name ILIKE $` + strOrder
			ParamsForGet = append(ParamsForGet, `%`+*Filter.Name+`%`)
			ParamsForCount = append(ParamsForCount, `%`+*Filter.Name+`%`)
			init++
		}

		if Filter.Genre != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.genre = $` + strOrder
			sqlCount += ` AND genre = $` + strOrder
			ParamsForGet = append(ParamsForGet, *Filter.Genre)
			ParamsForCount = append(ParamsForCount, *Filter.Genre)
			init++
		}

		if Filter.BeginDate != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.date >= $` + strOrder
			sqlCount += ` AND date >= $` + strOrder
			ParamsForGet = append(ParamsForGet, *Filter.BeginDate)
			ParamsForCount = append(ParamsForCount, *Filter.BeginDate)
			init++
		}

		if Filter.EndDate != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.date <= $` + strOrder
			sqlCount += ` AND date <= $` + strOrder
			ParamsForGet = append(ParamsForGet, *Filter.EndDate)
			ParamsForCount = append(ParamsForCount, *Filter.EndDate)
			init++
		}

		if Filter.Mechanics != nil {
			strOrder := strconv.Itoa(init)
			sqlGet += ` AND idea.mechanics @> $` + strOrder
			sqlCount += ` AND mechanics @> $` + strOrder
			ParamsForGet = append(ParamsForGet, *Filter.Mechanics)
			ParamsForCount = append(ParamsForCount, *Filter.Mechanics)
			init++
		}

		var SortingWayByIdeaLabel string
		SortingWayByIdeaLabel = "DESC"

		if Filter.Condition != nil {
			if *Filter.Condition == responses.RatedIdea {
				SortingWayByIdeaLabel = "ASC"
			}
		}
		strOrderLimit := strconv.Itoa(init)
		strOrderOffset := strconv.Itoa(init + 1)

		sqlGet += ` ORDER BY userrel.mark ` + SortingWayByIdeaLabel + `, idea.date DESC` + ` LIMIT $` + strOrderLimit + ` OFFSET $` + strOrderOffset
		ParamsForGet = append(ParamsForGet, Filter.Limit, Filter.Offset)
		if Filter.Offset == 0 {
			row := conn.QueryRow(ctx, sqlCount, ParamsForCount...)
			err = row.Scan(&item.Total)
			if err != nil {
				eMsg := "error in sqlCountIdea"
				clog.WithError(err).Error(eMsg)
				err = errors.Wrap(err, eMsg)
				return
			}

		}

		rows, err := conn.Query(ctx, sqlGet, ParamsForGet...)
		if err != nil {
			eMsg := "error in sqlGetIdeaList"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
		for rows.Next() {
			idea := models.IdeaLightData{}
			var mark *string
			err = rows.Scan(
				&idea.ID,
				&idea.Name,
				&idea.Date,
				&idea.Description,
				&idea.Worker.ID,
				&idea.Worker.Firstname,
				&idea.Worker.Lastname,
				&idea.Worker.Position,
				&idea.FilePath,
				&mark,
			)
			if err != nil {
				eMsg := "error ocurred while scanning sqlGetIdeaList"
				clog.WithError(err).Error(eMsg)
				err = errors.Wrap(err, eMsg)
				return
			}

			if mark != nil {
				idea.IsItNew = false
			} else {
				idea.IsItNew = true
			}

			item.Result = append(item.Result, idea)
		}

		numOfIdeas := len(item.Result)
		for i := 0; i < numOfIdeas; i++ {
			if !item.Result[i].IsItNew {
				row := conn.QueryRow(ctx, sqlGetOverAllRateIdea, item.Result[i].ID)
				var rateNum, rateSum int
				err = row.Scan(&rateNum, &rateSum)
				if err != nil {
					eMsg := "error in sqlGetOverAllRateIdea"
					clog.WithError(err).Error(eMsg)
					err = errors.Wrap(err, eMsg)
					return
				}

				if rateNum != 0 {
					point := float64(rateSum) / float64(rateNum)
					under := point - float64(int(point))
					upper := 1 - under
					if under >= upper {
						item.Result[i].OverallRate = int(point) + 1
					} else {
						item.Result[i].OverallRate = int(point)
					}

				}
			}
		}
		if Filter.WorkerID != nil {
			var lastSubmit *time.Time
			lastSubmitRow := conn.QueryRow(ctx, sqlSelectLastIdeaSubmittedDateWorker, *Filter.WorkerID)
			err = lastSubmitRow.Scan(&lastSubmit)
			if err != nil && err != pgx.ErrNoRows {
				eMsg := "error in sqlSelectLastIdeaSubmittedDateWorker"
				clog.WithError(err).Error(eMsg)
				err = errors.Wrap(err, eMsg)
				return
			}

			if err == nil {
				item.LastSubmitted = *lastSubmit
			}

			if err == pgx.ErrNoRows {
				err = nil
			}
		}
		return
	})
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *PgAccess) IdeaDelete(
	ctx context.Context,
	ID string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaDelete",
	})

	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {

		_, err = conn.Exec(ctx, sqlDeleteIdea, ID)
		if err != nil {
			eMsg := "error in sqlDeleteIdea"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		return
	})
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *PgAccess) IdeaUpdate(
	ctx context.Context,
	pTx pgx.Tx,
	NewIdea *models.IdeaUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaUpdate",
	})
	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {
		rollback = true
		_, err = tx.Exec(ctx,
			sqlUpdateIdea,
			NewIdea.Name,
			NewIdea.WorkerID,
			NewIdea.Date,
			NewIdea.Genre,
			NewIdea.Mechanics,
			NewIdea.Description,
			NewIdea.ID,
		)
		if err != nil {
			eMsg := "error in sqlUpdateIdea"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		_, err = tx.Exec(ctx, sqlDeleteIdeaLinks, NewIdea.ID)
		if err != nil {
			eMsg := "error in sqlDeleteIdeaLinks"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		linkNumber := len(NewIdea.Links)
		for i := 0; i < linkNumber; i++ {
			_, err = tx.Exec(ctx, sqlCreateLinkIdea, NewIdea.Links[i].Label, NewIdea.Links[i].URL, NewIdea.ID)
			if err != nil {
				eMsg := "error in sqlCreateLinkIdea"
				clog.WithError(err).Error(eMsg)
				err = errors.Wrap(err, eMsg)
				return
			}
		}

		return false, nil
	})
	if err != nil {
		eMsg := "error in d.runInTX()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

// GENRE

func (d *PgAccess) GenreUpdate(
	ctx context.Context,
	GenreUpdate models.GenreUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GenreUpdate",
	})

	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(ctx, sqlUpdateGenreName, GenreUpdate.NewGenre, GenreUpdate.OldGenre)
		if err != nil {
			eMsg := "error in sqlUpdateGenreName"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		_, err = conn.Exec(ctx, sqlUpdateAllGenreNamesIdea, GenreUpdate.NewGenre, GenreUpdate.OldGenre)
		if err != nil {
			eMsg := "error in sqlUpdateAllGenreNamesIdea"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
		return
	})
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

// MECHANICS

func (d *PgAccess) MechanicUpdate(
	ctx context.Context,
	MechUpdate models.MechanicUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.MechanicUpdate",
	})

	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(ctx, sqlUpdateMechanic, MechUpdate.NewMech, MechUpdate.OldMech)
		if err != nil {
			eMsg := "error in sqlUpdateMechanic"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		_, err = conn.Exec(ctx, sqlUpdateAllmechanicNamesIdea, MechUpdate.OldMech, MechUpdate.NewMech)
		if err != nil {
			eMsg := "error in sqlUpdateAllmechanicNamesIdea"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
		return
	})
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

// CRITERIA

func (d *PgAccess) CountCriteriaRates(
	ctx context.Context,
	ID string,
) (item int, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CountCriteriaRates",
	})

	err = d.runQuery(ctx, clog, func(conn *pgxpool.Conn) (err error) {

		var total int
		row := conn.QueryRow(ctx, sqlCountCriteriaRates, ID)
		err = row.Scan(&total)
		if err != nil {
			eMsg := "error in sqlCountCriteriaRates"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		item = total

		return
	})
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *PgAccess) CriteriaDelete(
	ctx context.Context,
	pTx pgx.Tx,
	ID string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaDelete",
	})

	err = d.runInTx(ctx, pTx, clog, func(tx pgx.Tx) (rollback bool, err error) {

		rollback = true

		_, err = tx.Exec(
			ctx,
			sqlDeleteCriteria,
			ID,
		)
		if err != nil {
			eMsg := "An error occurred on sqlDeleteCriteria"
			clog.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		return false, nil
	})

	if err != nil {
		eMsg := "An error occurred on d.runInTx()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////////MONGO///////////////////////////////////////////////////////////////////////

//IDEA

func (d *MgAccess) IdeaCreate(
	ctx context.Context,
	Idea *models.IdeaCreate,
) (err error) {

	_ = log.WithFields(log.Fields{
		"method": "PgAccess.IdeaCreate",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}

	db := client.Database("idea-share")
	coll := db.Collection("idea")

	rates := make([]models.RatingStructInIdea, 0)

	_, err = coll.InsertOne(ctx, bson.D{
		{Key: "name", Value: Idea.Name},
		{Key: "worker", Value: Idea.Worker},
		{Key: "date", Value: Idea.Date},
		{Key: "genre", Value: Idea.Genre},
		{Key: "mechanics", Value: Idea.Mechanics},
		{Key: "links", Value: Idea.Links},
		{Key: "description", Value: Idea.Description},
		{Key: "paths", Value: Idea.Paths},
		{Key: "rates", Value: rates},
	})

	return
}

func (d *MgAccess) IdeaRate(
	ctx context.Context,
	Rating *models.RateIdeaCritera,
) (item *int, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaRate",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}

	db := client.Database("idea-share")
	coll := db.Collection("idea")

	MatchStage := bson.M{"_id": Rating.IdeaID}
	removeOldRating := bson.M{"$pull": bson.M{"rates": bson.M{"user_id": Rating.Rating.UserID, "criteria_id": Rating.Rating.CriteriaID}}}
	addNewRating := bson.M{"$push": bson.M{"rates": Rating.Rating}}

	_, err = coll.UpdateOne(ctx, MatchStage, removeOldRating)
	if err != nil {
		eMsg := "error in removing old rating"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	_, err = coll.UpdateOne(ctx, MatchStage, addNewRating)
	if err != nil {
		eMsg := "error in addng new rate"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	options := options.FindOne()
	options.Projection = bson.M{"_id": 0, "rates": 1}

	var s models.ArrayOfRatesIdea

	_ = coll.FindOne(ctx, MatchStage, options).Decode(&s)
	if err != nil {
		eMsg := "Error in aggregation of CheckAllMechanicsArePresent"
		clog.WithError(err).Error(eMsg)
		return
	}
	var rateSum int
	rateNum := len(s.Rates)
	for i := 0; i < rateNum; i++ {
		rateSum += s.Rates[i].Rate
	}

	var averageRate int
	item = &averageRate

	if rateNum != 0 {
		point := float64(rateSum) / float64(rateNum)
		under := point - float64(int(point))
		upper := 1 - under
		if under >= upper {
			averageRate = int(point) + 1
		} else {
			averageRate = int(point)
		}

	}

	return
}

//GENRE
func (d *MgAccess) GenreUpsert(
	ctx context.Context,
	GenreName string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GenreUpsert",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("genre")

	filter := bson.M{"name": GenreName}
	update := bson.M{"$set": bson.M{"name": GenreName}}
	opts := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		eMsg := "error in Upserting genre"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}

func (d *MgAccess) GenreList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GenreList",
	})

	GENRES := make([]string, 0)
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
	db := client.Database("idea-share")
	coll := db.Collection("genre")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
	}
	var genres []bson.M

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &genres); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(genres); i++ {
		ps := genres[i]["name"].(string)
		GENRES = append(GENRES, ps)
	}

	item = &GENRES

	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *MgAccess) GenreDelete(
	ctx context.Context,
	GenreName string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GenreDelete",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("genre")

	_, err = coll.DeleteOne(ctx, bson.M{"name": GenreName})
	if err != nil {
		eMsg := "Error in genre delete"
		clog.WithError(err).Error(eMsg)
		return
	}
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

// MECHANICS

func (d *MgAccess) MechanicUpsert(
	ctx context.Context,
	Mechanics string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.MechanicUpsert",
	})
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("mechanic")

	filter := bson.M{"name": Mechanics}
	update := bson.M{"$set": bson.M{"name": Mechanics}}
	opts := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		eMsg := "error in Upserting Mechanics"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *MgAccess) MechanicList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.MechanicList",
	})
	MECHS := make([]string, 0)
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}

	db := client.Database("idea-share")
	coll := db.Collection("mechanic")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
	}
	var mechanics []bson.M

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &mechanics); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(mechanics); i++ {
		mech := mechanics[i]["name"].(string)
		MECHS = append(MECHS, mech)
	}

	item = &MECHS

	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}

	return
}

func (d *MgAccess) MechanicDelete(
	ctx context.Context,
	MechName string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.MechanicDelete",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("mechanic")

	_, err = coll.DeleteOne(ctx, bson.M{"name": MechName})
	if err != nil {
		eMsg := "Error in genre delete"
		clog.WithError(err).Error(eMsg)
		return
	}
	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}
	return
}

func (d *MgAccess) CheckAllMechanicsArePresent(
	ctx context.Context,
	mechList []string,
) (item bool, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CheckAllMechanicsArePresent",
	})
	item = false
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	groupStage := bson.D{{"$group", bson.D{{"_id", ""}, {"mechs", bson.D{{"$push", "$name"}}}}}}
	matchStage := bson.D{{"$match", bson.D{{"mechs", bson.D{{"$all", mechList}}}}}}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{groupStage, matchStage})
	if err != nil {
		eMsg := "Error in aggregation of CheckAllMechanicsArePresent"
		clog.WithError(err).Error(eMsg)
		return
	}

	var results []bson.M

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &results); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	if results != nil {
		item = true
	}

	return
}

//CRITERIA

func (d *MgAccess) CriteriaCreate(
	ctx context.Context,
	CriteriaName string,
) (item *models.CriteriaSpecData, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaCreate",
	})
	defer func() {
		if err != nil {
			item = nil
		}
	}()
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	row, err := coll.InsertOne(ctx, bson.D{
		{Key: "name", Value: CriteriaName},
	})
	if err != nil {
		eMsg := "An error occurred on Insert one"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &models.CriteriaSpecData{
		ID:   row.InsertedID.(primitive.ObjectID).Hex(),
		Name: CriteriaName,
	}
	return
}

func (d *MgAccess) CriteriaGetByName(
	ctx context.Context,
	criteriaName string,
) (item *string, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaGetByName",
	})
	defer func() {
		if err != nil {
			item = nil
		}
	}()
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	var u bson.M
	err = coll.FindOne(ctx, bson.M{"name": criteriaName}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			item = nil
			return
		}
		eMsg := "Error in Find criteria with name"
		clog.WithError(err).Error(eMsg)
		return
	}

	id := u["_id"].(primitive.ObjectID).Hex()
	item = &id

	return
}

func (d *MgAccess) CriteriaGetByID(
	ctx context.Context,
	ID primitive.ObjectID,
) (item *models.CriteriaSpecData, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaGetByID",
	})
	item = &models.CriteriaSpecData{}

	defer func() {
		if err != nil {
			item = nil
		}
	}()
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	var u bson.M
	err = coll.FindOne(ctx, bson.M{"_id": ID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			item = nil
			return
		}
		eMsg := "Error in Find criteria with ID"
		clog.WithError(err).Error(eMsg)
		return
	}

	item.ID = u["_id"].(primitive.ObjectID).Hex()
	item.Name = u["name"].(string)

	return
}

func (d *MgAccess) CriteriaUpdate(
	ctx context.Context,
	criter *models.CriteriaUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaUpdate",
	})
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	filter := bson.M{"_id": criter.ID}
	update := bson.M{"$set": bson.M{"name": criter.Name}}

	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		eMsg := "error in Update criteria"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	// TODO: idealaryn icinden hem update et

	if err != nil {
		eMsg := "Error in d.runQuery()"
		clog.WithError(err).Error(eMsg)
	}

	return
}

func (d *MgAccess) CriteriaList(
	ctx context.Context,
) (item *[]models.CriteriaSpecData, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaList",
	})

	criterias := make([]models.CriteriaSpecData, 0)
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
	db := client.Database("idea-share")
	coll := db.Collection("criteria")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
		return
	}
	var crtrs []bson.M

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &crtrs); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(crtrs); i++ {
		c := models.CriteriaSpecData{
			ID:   crtrs[i]["_id"].(primitive.ObjectID).Hex(),
			Name: crtrs[i]["name"].(string),
		}

		criterias = append(criterias, c)
	}

	item = &criterias
	return
}
