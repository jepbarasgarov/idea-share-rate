package datastore

import (
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///////////////////////////////////////////////////////////////////////////////////////////MONGO///////////////////////////////////////////////////////////////////////

//IDEA

func (d *MgAccess) IdeaCreate(
	ctx context.Context,
	Idea *models.IdeaCreate,
) (err error) {

	_ = log.WithFields(log.Fields{
		"method": "PgAccess.IdeaCreate",
	})

	db := d.client.Database("idea-share")
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
		{Key: "files", Value: Idea.AllFiles},
		{Key: "rates", Value: rates},
		{Key: "create_ts", Value: time.Now().UTC()},
		{Key: "update_ts", Value: time.Now().UTC()},
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

	db := d.client.Database("idea-share")
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
		eMsg := "error in adding new rate"
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

func (d *MgAccess) IdeaList(
	ctx context.Context,
	cu *responses.ActionInfo,
	Filter *models.IdeaFilter,
) (item *models.IdeaList, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaList",
	})

	item = &models.IdeaList{}
	item.Result = make([]models.IdeaLightData, 0)

	db := d.client.Database("idea-share")
	coll := db.Collection("idea")

	x := bson.D{}

	if Filter.WorkerID != nil {
		x = append(x, bson.E{"worker._id", *Filter.WorkerID})
	}

	if Filter.Name != nil {
		x = append(x, bson.E{"name", bson.D{{"$regex", *Filter.Name}, {"$options", "i"}}})

	}

	if Filter.Genre != nil {
		x = append(x, bson.E{"genre", *Filter.Genre})
	}

	if Filter.BeginDate != nil {
		dateCompareBegin := bson.E{"date", bson.M{"$gte": *Filter.BeginDate}}
		x = append(x, dateCompareBegin)
	}

	if Filter.EndDate != nil {
		dateCompareEnd := bson.E{"date", bson.M{"$lte": *Filter.EndDate}}
		x = append(x, dateCompareEnd)

	}

	if Filter.Mechanics != nil {
		if len(*Filter.Mechanics) > 0 {
			mechSearch := bson.E{"mechanics", bson.M{"$all": *Filter.Mechanics}}
			x = append(x, mechSearch)
		}
	}
	matchStageList := bson.D{{"$match", x}}
	unWindStagsList := bson.D{{"$unwind", bson.D{{"path", "$rates"}, {"preserveNullAndEmptyArrays", true}}}}
	projectStageLIst := bson.D{{"$project", bson.D{
		{"_id", 1},
		{"name", 1},
		{"worker", 1},
		{"date", 1},
		{"description", 1},
		{"create_ts", 1},
		{"path", bson.M{"$arrayElemAt": bson.A{"$files.file_path", 0}}},
		{"is_it_new", bson.M{"$ne": bson.A{"$rates.user_id", Filter.UserID}}},
		{"rate", "$rates.rate"},
	}}}
	groupStageList := bson.D{{"$group", bson.D{
		{"_id", "$_id"},
		{"name", bson.M{"$first": "$name"}},
		{"worker", bson.M{"$first": "$worker"}},
		{"date", bson.M{"$first": "$date"}},
		{"create_ts", bson.M{"$first": "$create_ts"}},
		{"description", bson.M{"$first": "$description"}},
		{"path", bson.M{"$first": "$path"}},
		{"is_it_new", bson.M{"$min": "$is_it_new"}},
		{"avg", bson.M{"$avg": "$rate"}},
	}}}

	var sortStageList bson.D
	sortStageList = bson.D{{"$sort", bson.D{{"is_it_new", -1}, {"create_ts", -1}}}}

	if Filter.Condition != nil {
		if *Filter.Condition == responses.RatedIdea {
			sortStageList = bson.D{{"$sort", bson.M{"is_it_new": 1}}}
		}
	}

	limitStageList := bson.D{{"$limit", Filter.Limit}}
	offsetStageList := bson.D{{"$skip", Filter.Offset}}

	cursorIdeaLits, err := coll.Aggregate(ctx, mongo.Pipeline{matchStageList, unWindStagsList, projectStageLIst, groupStageList, sortStageList, offsetStageList, limitStageList})
	if err != nil {
		eMsg := "Error in cursorIdeaLits"
		clog.WithError(err).Error(eMsg)
		return
	}

	for cursorIdeaLits.Next(ctx) {
		if err = cursorIdeaLits.All(ctx, &item.Result); err != nil {
			eMsg := "Error in reading cursorIdeaLits"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(item.Result); i++ {
		if !item.Result[i].IsItNew {
			under := *item.Result[i].AvgRate - float64(int(*item.Result[i].AvgRate))
			upper := 1 - under
			if under >= upper {
				item.Result[i].OverallRate = int(*item.Result[i].AvgRate) + 1
			} else {
				item.Result[i].OverallRate = int(*item.Result[i].AvgRate)
			}

			item.Result[i].AvgRate = nil
		}

		item.Result[i].AvgRate = nil

	}

	groupStageCount := bson.D{{"$group", bson.D{
		{"_id", ""},
		{"total", bson.M{"$sum": 1}},
	}}}

	projectStageCount := bson.D{{"$project", bson.D{{"_id", 0}}}}

	var total []bson.M

	cursorIdeaCount, err := coll.Aggregate(ctx, mongo.Pipeline{matchStageList, groupStageCount, projectStageCount})
	if err != nil {
		eMsg := "Error in cursorIdeaCount"
		clog.WithError(err).Error(eMsg)
		return
	}

	for cursorIdeaCount.Next(ctx) {
		if err = cursorIdeaCount.All(ctx, &total); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	if len(total) > 0 {
		item.Total = int(total[0]["total"].(int32))
	}

	if Filter.WorkerID != nil {
		matchStageSubmit := bson.D{{"$match", bson.M{"worker._id": *Filter.WorkerID}}}
		sortStageSubmit := bson.D{{"$sort", bson.M{"create_ts": -1}}}
		limitStageSubmit := bson.D{{"$limit", 1}}
		projectStageSubmit := bson.D{{"$project", bson.D{{"create_ts", 1}, {"_id", 0}}}}
		cursorGetLastSubmitOfWorker, err1 := coll.Aggregate(ctx, mongo.Pipeline{matchStageSubmit, sortStageSubmit, limitStageSubmit, projectStageSubmit})
		if err1 != nil {
			eMsg := "Error in cursorGetLastSubmitOfWorker"
			clog.WithError(err1).Error(eMsg)
			return
		}

		var lastSubmit []bson.M

		for cursorGetLastSubmitOfWorker.Next(ctx) {
			if err = cursorGetLastSubmitOfWorker.All(ctx, &lastSubmit); err != nil {
				eMsg := "Error in reading cursorGetLastSubmitOfWorker"
				clog.WithError(err).Error(eMsg)
				return
			}
		}

		if len(lastSubmit) > 0 {
			item.LastSubmitted = lastSubmit[0]["create_ts"].(primitive.DateTime).Time().UTC()

		}

	}

	return
}

func (d *MgAccess) IdeaGet(
	ctx context.Context,
	cu *responses.ActionInfo,
	ID primitive.ObjectID,
) (item *models.IdeaSpecData, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GetIdeaByID",
	})
	item = &models.IdeaSpecData{}

	db := d.client.Database("idea-share")
	collIdea := db.Collection("idea")

	err = collIdea.FindOne(ctx, bson.M{"_id": ID}).Decode(&item)
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

	usersOwnRates := make([]responses.CriteriaRate, 0)
	totalRate := 0
	rateNum := 0
	isItNewToUser := true
	for i := 0; i < len(item.CriteriaRates); i++ {
		totalRate += item.CriteriaRates[i].Rate
		rateNum++
		if *item.CriteriaRates[i].UserID == cu.ID {
			item.CriteriaRates[i].UserID = nil
			usersOwnRates = append(usersOwnRates, item.CriteriaRates[i])
			isItNewToUser = false
		}
	}

	item.CriteriaRates = usersOwnRates

	collCriteria := db.Collection("criteria")

	cursor, err := collCriteria.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
		return
	}
	var crtrs []models.CriteriaSpecData

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &crtrs); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(crtrs); i++ {
		unequal := 0
		for j := 0; j < len(item.CriteriaRates); j++ {
			if crtrs[i].ID != item.CriteriaRates[j].ID {
				unequal++
			} else {
				item.CriteriaRates[j].Name = crtrs[i].Name
				break
			}
		}
		if unequal == len(item.CriteriaRates) {
			criterRate := responses.CriteriaRate{
				ID:   crtrs[i].ID,
				Name: crtrs[i].Name,
				Rate: 0,
			}

			item.CriteriaRates = append(item.CriteriaRates, criterRate)
		}
	}

	if !isItNewToUser {
		avg := float64(totalRate) / float64(rateNum)
		under := avg - float64(int(avg))
		upper := 1 - under
		if under >= upper {
			item.OverallRate = int(avg) + 1
		} else {
			item.OverallRate = int(avg)
		}

	}

	return
}

func (d *MgAccess) IdeaDelete(
	ctx context.Context,
	ID primitive.ObjectID,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaDelete",
	})

	db := d.client.Database("idea-share")
	collIdea := db.Collection("idea")

	_, err = collIdea.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		eMsg := "Error in Idea delete"
		clog.WithError(err).Error(eMsg)
		return
	}

	return
}

func (d *MgAccess) IdeaUpdate(
	ctx context.Context,
	pTx pgx.Tx,
	NewIdea *models.IdeaUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.IdeaUpdate",
	})

	db := d.client.Database("idea-share")
	collIdea := db.Collection("idea")

	filter := bson.M{"_id": NewIdea.ID}
	update := bson.M{"$set": bson.M{
		"name":        NewIdea.Name,
		"date":        NewIdea.Date,
		"worker":      NewIdea.Worker,
		"description": NewIdea.Description,
		"genre":       NewIdea.Genre,
		"mechanics":   NewIdea.Mechanics,
		"links":       NewIdea.Links,
		"update_ts":   time.Now().UTC(),
	}}

	_, err = collIdea.UpdateOne(ctx, filter, update)
	if err != nil {
		eMsg := "error in Updating idea"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
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

	db := d.client.Database("idea-share")
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

	db := d.client.Database("idea-share")
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

	db := d.client.Database("idea-share")
	coll := db.Collection("genre")

	_, err = coll.DeleteOne(ctx, bson.M{"name": GenreName})
	if err != nil {
		eMsg := "Error in genre delete"
		clog.WithError(err).Error(eMsg)
		return
	}

	return
}

func (d *MgAccess) GenreUpdate(
	ctx context.Context,
	GenreUpdate models.GenreUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.GenreUpdate",
	})

	db := d.client.Database("idea-share")
	collGenre := db.Collection("genre")
	collIdea := db.Collection("idea")

	filterGenre := bson.M{"name": GenreUpdate.OldGenre}
	updateGenre := bson.M{"$set": bson.M{"name": GenreUpdate.NewGenre}}

	_, err = collGenre.UpdateOne(ctx, filterGenre, updateGenre)
	if err != nil {
		eMsg := "Error in genre update"
		clog.WithError(err).Error(eMsg)
		return
	}

	filterIdea := bson.M{"genre": GenreUpdate.OldGenre}
	updateIdea := bson.M{"$set": bson.M{"genre": GenreUpdate.NewGenre}}

	_, err = collIdea.UpdateMany(ctx, filterIdea, updateIdea)
	if err != nil {
		eMsg := "Error in genre update from idea"
		clog.WithError(err).Error(eMsg)
		return
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

	db := d.client.Database("idea-share")
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

	db := d.client.Database("idea-share")
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

	db := d.client.Database("idea-share")
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

	db := d.client.Database("idea-share")
	coll := db.Collection("mechanic")

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

func (d *MgAccess) MechanicUpdate(
	ctx context.Context,
	MechUpdate models.MechanicUpdate,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.MechanicUpdate",
	})

	db := d.client.Database("idea-share")
	collMech := db.Collection("mechanic")
	collIdea := db.Collection("idea")

	filterMech := bson.M{"name": MechUpdate.OldMech}
	updateMech := bson.M{"$set": bson.M{"name": MechUpdate.NewMech}}

	_, err = collMech.UpdateOne(ctx, filterMech, updateMech)
	if err != nil {
		eMsg := "Error in genre update"
		clog.WithError(err).Error(eMsg)
		return
	}

	MatchStage := bson.M{"mechanics": MechUpdate.OldMech}
	removeOldMech := bson.M{"$pull": bson.M{"mechanics": MechUpdate.OldMech}}
	addNewMech := bson.M{"$push": bson.M{"mechanics": MechUpdate.NewMech}}

	_, err = collIdea.UpdateMany(ctx, MatchStage, addNewMech)
	if err != nil {
		eMsg := "error in adding new mechanic"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	_, err = collIdea.UpdateMany(ctx, MatchStage, removeOldMech)
	if err != nil {
		eMsg := "error in removing old mechanic"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
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

	db := d.client.Database("idea-share")
	coll := db.Collection("criteria")

	row, err := coll.InsertOne(ctx, bson.D{
		{Key: "name", Value: CriteriaName},
		{Key: "create_ts", Value: time.Now().UTC()},
		{Key: "update_ts", Value: time.Now().UTC()},
	})
	if err != nil {
		eMsg := "An error occurred on Insert one"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &models.CriteriaSpecData{
		ID:   row.InsertedID.(primitive.ObjectID),
		Name: CriteriaName,
	}

	return
}

func (d *MgAccess) CriteriaGetByName(
	ctx context.Context,
	criteriaName string,
) (item *primitive.ObjectID, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaGetByName",
	})
	defer func() {
		if err != nil {
			item = nil
		}
	}()

	db := d.client.Database("idea-share")
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

	id := u["_id"].(primitive.ObjectID)
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

	db := d.client.Database("idea-share")
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

	item.ID = u["_id"].(primitive.ObjectID)
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

	db := d.client.Database("idea-share")
	collCriteria := db.Collection("criteria")

	filter := bson.M{"_id": criter.ID}
	update := bson.M{"$set": bson.M{"name": criter.Name, "update_ts": time.Now().UTC()}}

	_, err = collCriteria.UpdateOne(ctx, filter, update)
	if err != nil {
		eMsg := "error in Update criteria"
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

func (d *MgAccess) CriteriaList(
	ctx context.Context,
) (item *[]models.CriteriaSpecData, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaList",
	})

	criterias := make([]models.CriteriaSpecData, 0)

	db := d.client.Database("idea-share")
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
			ID:   crtrs[i]["_id"].(primitive.ObjectID),
			Name: crtrs[i]["name"].(string),
		}

		criterias = append(criterias, c)
	}

	item = &criterias
	return
}

func (d *MgAccess) CriteriaDelete(
	ctx context.Context,
	pTx pgx.Tx,
	ID primitive.ObjectID,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CriteriaDelete",
	})

	db := d.client.Database("idea-share")
	coll := db.Collection("criteria")

	_, err = coll.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		eMsg := "error in delete criteria"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}

func (d *MgAccess) CountCriteriaRates(
	ctx context.Context,
	ID primitive.ObjectID,
) (item int, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CountCriteriaRates",
	})

	//TODO : su functiony goni idelan nacesinde bar dp sanap optimise etjek bol

	db := d.client.Database("idea-share")
	collIdea := db.Collection("idea")

	unWindStagsList := bson.D{{"$unwind", bson.D{{"path", "$rates"}, {"preserveNullAndEmptyArrays", false}}}}
	projectStageLIst := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"rates.criteria_id", 1},
	}}}
	groupStageList := bson.D{{"$group", bson.D{
		{"_id", "$rates.criteria_id"},
		{"count", bson.M{"$sum": 1}},
	}}}
	matchStageList := bson.D{{"$match", bson.M{"_id": ID}}}

	cursorCountCriteria, err := collIdea.Aggregate(ctx, mongo.Pipeline{projectStageLIst, unWindStagsList, groupStageList, matchStageList})
	if err != nil {
		eMsg := "Error in cursorCountCriteria"
		clog.WithError(err).Error(eMsg)
		return
	}
	var res []bson.M
	for cursorCountCriteria.Next(ctx) {
		if err = cursorCountCriteria.All(ctx, &res); err != nil {
			eMsg := "Error in reading cursorCountCriteria"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	if len(res) > 0 {
		item = int(res[0]["count"].(int32))
	}

	return
}
