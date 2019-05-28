package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"time"
)

var client *mongo.Client
var databases = make(map[string]*Database)

type Database struct {
	name string
}

func GetDatabase(name string) (database *Database) {
	database = databases[name]
	if database == nil {
		database = NewDatabase(name)
	}
	return
}

func NewDatabase(database string) *Database {
	return &Database{name: database}
}

func checkClient() error {
	if client == nil || client.Ping(context.TODO(), nil) != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err == nil {
			err = client.Ping(context.TODO(), nil)
			if err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

func (d *Database) collection(collection string) (*mongo.Collection, error) {
	if err := checkClient(); err != nil {
		log.Println(err)
		return nil, err
	}
	return client.Database(d.name).Collection(collection), nil
}

func (d *Database) FindOne(collection string, query interface{}, options *options.FindOneOptions, ret interface{}) error {
	if c, e := d.collection(collection); e == nil {
		err := c.FindOne(context.TODO(), query, options).Decode(ret)
		if err != nil {
			return err
		}
		return nil
	} else {
		return e
	}
}

func (d *Database) Count(collection string, query interface{}, options *options.CountOptions) (int64, error) {
	if c, e := d.collection(collection); e == nil {
		return c.CountDocuments(context.TODO(), query, options)
	} else {
		return 0, e
	}
}

func (d *Database) FindMany(collection string, query interface{}, options *options.FindOptions, ret interface{}) error {
	if c, e := d.collection(collection); e == nil {
		cur, err := c.Find(context.TODO(), query, options)
		if err != nil {
			return err
		}
		typ := reflect.TypeOf(ret)
		if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Slice {
			return errors.New("ret must be a pointer of slice")
		}
		ptr := reflect.ValueOf(ret)
		set := ptr.Elem()
		//set = set.Slice(0, set.Cap())
		for cur.Next(context.TODO()) {
			ele := reflect.New(typ.Elem().Elem())
			err := cur.Decode(ele.Interface())
			if err != nil {
				log.Fatalln(err)
			}
			set = reflect.Append(set, ele.Elem())
			//set = set.Slice(0, set.Cap())
		}
		if err := cur.Close(context.TODO()); err != nil {
			log.Fatalln(err)
		}
		ptr.Elem().Set(set)
		return nil
	} else {
		return e
	}
}

func (d *Database) InsertOne(collection string, v interface{}) (*mongo.InsertOneResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.InsertOne(context.TODO(), v)
		return ret, err
	} else {
		return nil, e
	}
}

func (d *Database) InsertMany(collection string, vs []interface{}, options *options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.InsertMany(context.TODO(), vs, options)
		return ret, err
	} else {
		return nil, e
	}
}

func (d *Database) UpdateOne(collection string, query interface{}, update interface{}, options *options.UpdateOptions) (*mongo.UpdateResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.UpdateOne(context.TODO(), query, update, options)
		return ret, err
	} else {
		return nil, e
	}
}

func (d *Database) UpdateMany(collection string, query interface{}, update interface{}, options *options.UpdateOptions) (*mongo.UpdateResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.UpdateMany(context.TODO(), query, update, options)
		return ret, err
	} else {
		return nil, e
	}
}

func (d *Database) DeleteOne(collection string, query interface{}) (*mongo.DeleteResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.DeleteOne(context.TODO(), query)
		return ret, err
	} else {
		return nil, e
	}
}

func (d *Database) DeleteMany(collection string, query interface{}) (*mongo.DeleteResult, error) {
	if c, e := d.collection(collection); e == nil {
		ret, err := c.DeleteMany(context.TODO(), query)
		return ret, err
	} else {
		return nil, e
	}
}
