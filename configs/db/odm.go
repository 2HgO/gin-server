package db

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/2HgO/gin-web-server/errors"
)

type Collection struct {
	// collection is a connection to a MongoDB collection
	collection *mongo.Collection

	// database is a connection to the MongoDB database that houses the collection
	database *mongo.Database
}

func newCollection(d *mongo.Database, collection string) *Collection {
	if d == nil {
		panic("Invalid DB connection")
	}
	return &Collection{
		collection: d.Collection(collection),
		database:   d,
	}
}

type Doc interface {
	SetID(interface{})
	GetID() *primitive.ObjectID
	SetCreated(time.Time)
	SetUpdated(time.Time)
}

type SoftDoc interface {
	Doc
	SoftDelete()
}

type AfterDeleteHook interface {
	AfterDelete(*mongo.Database)
}

type AfterRecoverHook interface {
	AfterRecover(*mongo.Database)
}

type ValidationHook interface {
	Validate(*mongo.Database) error
}

type PreSaveHook interface {
	PreSave(*Collection) error
}

func (c *Collection) Database() *mongo.Database {
	return c.database
}

func (c *Collection) Collection() *mongo.Collection {
	return c.collection
}

func (c *Collection) Exists(query interface{}) (bool, error) {
	con, err := c.collection.Clone()
	if err != nil {
		return true, err // Return true so write operation does not create potential duplicate when cloning fails
	}
	iter, err := con.Find(nil, query, options.Find().SetLimit(1))
	return iter.TryNext(nil), err
}

func (c *Collection) DeletedExists(query interface{}) (bool, error) {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")
	return (&Collection{collection: deleted}).Exists(query)
}

func (c *Collection) Insert(doc Doc) error {
	if reflect.TypeOf(doc).Kind() != reflect.Ptr || doc == nil {
		return errors.ErrInvalidDoc
	}

	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	if d, ok := doc.(ValidationHook); ok {
		if err := d.Validate(con.Database()); err != nil {
			return err
		}
	}

	if d, ok := doc.(PreSaveHook); ok {
		if err := d.PreSave(c); err != nil {
			return err
		}
	}

	now := time.Now()
	if doc.GetID() == nil {
		doc.SetCreated(now)
	}
	doc.SetUpdated(now)

	res, err := con.InsertOne(nil, doc)
	if err != nil {
		return err
	}
	doc.SetID(res.InsertedID)

	return nil
}

func (c *Collection) FindOne(query interface{}, out Doc) error {
	if out == nil || reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.ErrInvalidDoc
	}

	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	if err := con.FindOne(nil, query).Decode(out); errors.Is(err, mongo.ErrNoDocuments) {
		return errors.NotFoundError
	} else {
		return err
	}
}

func (c *Collection) FindByID(id *primitive.ObjectID, out Doc) error {
	if out == nil || reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.ErrInvalidDoc
	}

	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	if err := con.FindOne(nil, bson.D{{Key: "_id", Value: id}}).Decode(out); errors.Is(err, mongo.ErrNoDocuments) {
		return errors.NotFoundError
	} else {
		return err
	}
}

func (c *Collection) Find(query interface{}) (*mongo.Cursor, error) {
	con, err := c.collection.Clone()
	if err != nil {
		return nil, err
	}

	return con.Find(nil, query)
}

func (c *Collection) UpdateOne(query interface{}, doc Doc) error {
	if doc == nil || reflect.TypeOf(doc).Kind() != reflect.Ptr {
		return errors.ErrInvalidDoc
	}

	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	if d, ok := doc.(ValidationHook); ok {
		if err := d.Validate(con.Database()); err != nil {
			return err
		}
	}

	if d, ok := doc.(PreSaveHook); ok {
		if err := d.PreSave(c); err != nil {
			return err
		}
	}

	doc.SetUpdated(time.Now())

	_, err = con.UpdateOne(nil, query, doc, options.Update().SetUpsert(true))
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.NotFoundError
	}
	return err
}

func (c *Collection) PartialUpdateOne(query, update interface{}, out Doc) error {
	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	res := con.FindOneAndUpdate(nil, query, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if out == nil || reflect.TypeOf(out).Kind() != reflect.Ptr {
		return res.Err()
	}
	if err := res.Decode(out); errors.Is(err, mongo.ErrNoDocuments) {
		return errors.NotFoundError
	} else {
		return err
	}
}

func (c *Collection) Update(query interface{}, update []Doc) (int64, error) {
	con, err := c.collection.Clone()
	if err != nil {
		return 0, err
	}

	now := time.Now()
	for _, val := range update {
		if reflect.TypeOf(val).Kind() != reflect.Ptr || val == nil {
			return 0, fmt.Errorf("Error in update slice: %w", errors.ErrInvalidDoc)
		}
		val.SetUpdated(now)
	}

	res, err := con.UpdateMany(nil, query, update, options.Update().SetUpsert(true))

	return res.ModifiedCount, err
}

func (c *Collection) DeleteOne(query interface{}, out Doc) error {
	con, err := c.collection.Clone()
	if err != nil {
		return err
	}

	res := con.FindOneAndDelete(nil, query)

	if reflect.TypeOf(out).Kind() != reflect.Ptr || out == nil {
		return res.Err()
	}

	err = res.Decode(out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.NotFoundError
		}
		return err
	}

	if d, ok := out.(AfterDeleteHook); ok {
		go d.AfterDelete(con.Database())
	}

	return nil
}

func (c *Collection) SoftDelete(query interface{}, out SoftDoc) error {
	if err := c.DeleteOne(query, out); err != nil {
		return err
	}

	if d, ok := out.(AfterDeleteHook); ok {
		go d.AfterDelete(c.Database())
	}

	deleted := c.database.Collection(c.collection.Name() + "_deleted")

	_, err := deleted.InsertOne(nil, out)
	return err
}

func (c *Collection) SoftRecover(query interface{}, out SoftDoc) error {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")

	res := deleted.FindOneAndDelete(nil, query)
	if err := res.Decode(out); err != nil {
		return err
	}

	out.SetUpdated(time.Now())
	err := c.Insert(out)

	if d, ok := out.(AfterRecoverHook); err == nil && ok {
		go d.AfterRecover(deleted.Database())
	}

	return err
}

func (c *Collection) FindOneDeleted(query interface{}, out SoftDoc) error {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")
	return (&Collection{collection: deleted}).FindOne(query, out)
}

func (c *Collection) FindDeletedByID(id *primitive.ObjectID, out SoftDoc) error {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")
	return (&Collection{collection: deleted}).FindByID(id, out)
}

func (c *Collection) FindDeleted(query interface{}) (*mongo.Cursor, error) {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")
	return (&Collection{collection: deleted}).Find(query)
}

func (c *Collection) Aggregate(pipeline interface{}) (*mongo.Cursor, error) {
	con, err := c.collection.Clone()
	if err != nil {
		return nil, err
	}

	return con.Aggregate(nil, pipeline, options.Aggregate().SetAllowDiskUse(true))
}

func (c *Collection) Count(query interface{}) (int64, error) {
	con, err := c.collection.Clone()
	if err != nil {
		return 0, err
	}

	return con.CountDocuments(nil, query)
}

func (c *Collection) CountDeleted(query interface{}) (int64, error) {
	deleted := c.database.Collection(c.collection.Name() + "_deleted")
	return (&Collection{collection: deleted}).Count(query)
}
