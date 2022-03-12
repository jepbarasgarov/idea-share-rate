package datastore

import (
	"belli/onki-game-ideas-mongo-backend/models"
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//////////////////////////////////////////////////////MONGO///////////////////////////////////////////////////////////////////////////////////////////

func (d *MgAccess) WorkerCreate(
	ctx context.Context,
	worker *models.WorkerCreate,
) (item *models.WorkerBsonModelInIdea, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.WorkerCreate",
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
	coll := db.Collection("worker")

	row, err := coll.InsertOne(ctx, bson.D{
		{Key: "firstname", Value: worker.Firstname},
		{Key: "lastname", Value: worker.Lastname},
		{Key: "position", Value: worker.Position},
	})
	if err != nil {
		eMsg := "An error occurred on Insert one"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &models.WorkerBsonModelInIdea{
		ID:        row.InsertedID.(primitive.ObjectID),
		Firstname: worker.Firstname,
		LastName:  worker.Lastname,
		Position:  worker.Position,
	}

	return
}

func (d *MgAccess) Workerget(
	ctx context.Context,
	ID primitive.ObjectID,
) (item *models.WorkerBsonModelInIdea, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.Workerget",
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
	coll := db.Collection("worker")

	var u models.WorkerBsonModelInIdea
	err = coll.FindOne(ctx, bson.M{"_id": ID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			item = nil
			return
		}
		eMsg := "Error in Find worker with ID"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &u

	return
}

func (d *MgAccess) WorkerDelete(
	ctx context.Context,
	ID primitive.ObjectID,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.WorkerDelete",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return

	}
	db := client.Database("idea-share")
	coll := db.Collection("worker")

	_, err = coll.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		eMsg := "Error in Find worker with ID"
		clog.WithError(err).Error(eMsg)
		return
	}

	return
}

func (d *MgAccess) WorkerUpdate(
	ctx context.Context,
	worker *models.WorkerUpdate,
) (item *models.WorkerBsonModelInIdea, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.WorkerUpdate",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := client.Database("idea-share")
	workerColl := db.Collection("worker")
	ideaColl := db.Collection("idea")

	filterWorker := bson.M{"_id": worker.ID}
	filterWorkerInIdea := bson.M{"worker._id": worker.ID}

	updateWorker := bson.M{"$set": bson.M{"firstname": worker.Firstname, "lastname": worker.Lastname, "position": worker.Position}}
	updateWorkerInIdea := bson.M{"$set": bson.M{"worker.firstname": worker.Firstname, "worker.lastname": worker.Lastname, "worker.position": worker.Position}}

	_, err = workerColl.UpdateOne(ctx, filterWorker, updateWorker)
	if err != nil {
		eMsg := "error in Updating worker"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	_, err = ideaColl.UpdateMany(ctx, filterWorkerInIdea, updateWorkerInIdea)
	if err != nil {
		eMsg := "error in Updating worker from related idea's"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}
	item = &models.WorkerBsonModelInIdea{
		ID:        worker.ID,
		Firstname: worker.Firstname,
		LastName:  worker.Lastname,
		Position:  worker.Position,
	}

	return

}

func (d *MgAccess) WorkerAutocompleteList(
	ctx context.Context,
) (item *[]models.WorkerBsonModelInIdea, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.WorkerAutocompleteList",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
	db := client.Database("idea-share")
	coll := db.Collection("worker")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
		return
	}
	var wrkrs []models.WorkerBsonModelInIdea

	if err = cursor.All(ctx, &wrkrs); err != nil {
		eMsg := "Error in reading cursor"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &wrkrs

	return
}

func (d *MgAccess) CountWorkersIdea(
	ctx context.Context,
	ID primitive.ObjectID,
) (item int, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.CountWorkersIdea",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}

	db := client.Database("idea-share")
	collIdea := db.Collection("idea")

	projectStageLIst := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"worker._id", 1},
	}}}
	groupStageList := bson.D{{"$group", bson.D{
		{"_id", "$worker._id"},
		{"count", bson.M{"$sum": 1}},
	}}}
	matchStageList := bson.D{{"$match", bson.M{"_id": ID}}}

	cursorCountIdeaOfWorker, err := collIdea.Aggregate(ctx, mongo.Pipeline{projectStageLIst, groupStageList, matchStageList})
	if err != nil {
		eMsg := "Error in cursorCountIdeaOfWorker"
		clog.WithError(err).Error(eMsg)
		return
	}
	var res []bson.M
	for cursorCountIdeaOfWorker.Next(ctx) {
		if err = cursorCountIdeaOfWorker.All(ctx, &res); err != nil {
			eMsg := "Error in reading cursorCountIdeaOfWorker"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	if len(res) > 0 {
		item = int(res[0]["count"].(int32))
	}

	return

}

// position

func (d *MgAccess) PositionUpsert(
	ctx context.Context,
	PositionName string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.PositionUpsert",
	})

	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
	db := client.Database("idea-share")
	coll := db.Collection("position")

	filter := bson.M{"name": PositionName}
	update := bson.M{"$set": bson.M{"name": PositionName}}
	opts := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		eMsg := "error in Upserting position"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}

func (d *MgAccess) PositionList(
	ctx context.Context,
) (item *[]string, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.PositionList",
	})

	positions := make([]string, 0)
	client, err := mongo.Connect(ctx, d.ClientOptions)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
	db := client.Database("idea-share")
	coll := db.Collection("position")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		eMsg := "Error in Find"
		clog.WithError(err).Error(eMsg)
	}
	var pstns []bson.M

	for cursor.Next(ctx) {
		if err = cursor.All(ctx, &pstns); err != nil {
			eMsg := "Error in reading cursor"
			clog.WithError(err).Error(eMsg)
			return
		}
	}

	for i := 0; i < len(pstns); i++ {
		ps := pstns[i]["name"].(string)
		positions = append(positions, ps)
	}

	item = &positions

	return

}
