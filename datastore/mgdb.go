package datastore

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PgAccess struct {
	Access,
	pool *pgxpool.Pool
}

type pgxWithTx func(tx pgx.Tx) (rollback bool, err error)
type pgxQuery func(conn *pgxpool.Conn) (err error)

func NewPgAccess(conf *config.Config) (pg *PgAccess, err error) {
	var pool *pgxpool.Pool
	pool, err = pgxpool.Connect(context.Background(), config.Conf.DbConn)
	if err != nil {
		eMsg := "error creating connection pool"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}
	pg = &PgAccess{pool: pool}
	return
}

func NewPgAccessWithPool(pool *pgxpool.Pool) (pg *PgAccess) {
	return &PgAccess{pool: pool}
}

func (d *PgAccess) runInTx(ctx context.Context, pTx pgx.Tx, clog *log.Entry, f pgxWithTx) (err error) {
	var conn *pgxpool.Conn
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	rollback := true
	var tx pgx.Tx
	defer func() {
		if rollback && tx != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil && err != pgx.ErrTxClosed {
				clog.WithError(rErr).Error("Error in tx.Rollback")
			}
		}
	}()
	if pTx == nil {
		conn, err = d.pool.Acquire(ctx)
		if err != nil {
			clog.WithError(err).Error("error acquiring connection")
			return
		}
		tx, err = conn.Begin(ctx)
		if err != nil {
			eMsg := "Error in conn.Begin"
			clog.WithError(err).Error(eMsg)
			return errors.Wrap(err, eMsg)
		}
	} else {
		tx, err = pTx.Begin(ctx)
		if err != nil {
			eMsg := "Error in tx.Begin"
			clog.WithError(err).Error(eMsg)
			return errors.Wrap(err, eMsg)
		}
	}
	rollback, err = f(tx)
	if err != nil {
		eMsg := "error in executing f"
		clog.WithError(err).Error(eMsg)
		return errors.Wrap(err, eMsg)
	}
	if !rollback {
		err = tx.Commit(ctx)
		if err != nil {
			eMsg := "Error in tx.Commit"
			clog.WithError(err).Error(eMsg)
			return errors.Wrap(err, eMsg)
		}
	}
	return
}

func (d *PgAccess) runQuery(ctx context.Context, clog *log.Entry, f pgxQuery) (err error) {
	var conn *pgxpool.Conn
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	conn, err = d.pool.Acquire(ctx)
	if err != nil {
		clog.WithError(err).Error("error acquiring connection")
		return
	}

	err = f(conn)
	if err != nil {
		eMsg := "error in executing f"
		clog.WithError(err).Error(eMsg)
		return errors.Wrap(err, eMsg)
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

type MgAccess struct {
	Access
	ClientOptions *options.ClientOptions
	client        *mongo.Client
}

func NewMgAccess(conf *config.Config) (mg *MgAccess, err error) {
	clientOptions := options.Client().ApplyURI(conf.MgConn)
	if err != nil {
		eMsg := "error creating connection pool"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	ValidateDatabase(client)

	mg = &MgAccess{ClientOptions: clientOptions, client: client}
	return

}

func ValidateDatabase(client *mongo.Client) {

	db := client.Database("idea-share")

	x := db.RunCommand(context.TODO(), bson.D{
		{"collMod", "worker"},
		{"validator", bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"firstname", "lastname", "position"},
				"properties": bson.M{
					"firstname": bson.M{
						"bsonType":    "string",
						"maxLength":   64,
						"minLength":   0,
						"description": "must be a string and is required",
					},
					"lastname": bson.M{
						"bsonType":    "string",
						"maxLength":   64,
						"minLength":   0,
						"description": "must be a string and is required",
					},
					"position": bson.M{
						"bsonType":    "string",
						"maxLength":   128,
						"minLength":   0,
						"description": "must be a string and is required",
					},
				},
			},
		}},
		{"validationLevel", "strict"},
	})

	fmt.Println(x.Err())
}
