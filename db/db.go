package db

import (
	"context"
	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"modoowiki/model"
	"modoowiki/model/config"
	"time"
)

type db struct {
	config.DB
	*mongo.Database
	page *mongo.Collection
	end  chan struct{}
	*golog.Logger
}

func New(conf config.DB, level golog.Level) (*db, error) {

	d := new(db)
	d.DB = conf

	d.Logger = golog.New()
	d.Logger.Level = level

	d.end = make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.DB.Path()))
	if err == nil {
		d.Database = client.Database(d.DB.Name)
		d.page = d.Collection(d.PageCollection.Name)
		d.Info("Finish initializing DB")
	}

	go d.receiveEnd(ctx, client)
	return d, err

}

func (d *db) receiveEnd(ctx context.Context, client *mongo.Client) {
	<-d.end
	if err := client.Disconnect(ctx); err != nil {
		d.Fatal(err)
	}
}

func (d *db) Close() {
	d.Info("Close Database connection")
	close(d.end)
}

func (d *db) GetPage(title string) (page model.Page, err error) {

	cursor, err := d.page.Find(context.TODO(), bson.D{{d.PageCollection.Title, title}})

	if err == nil && cursor.Next(context.TODO()) {
		page.Key = cursor.Current.Index(0).Value().String()
		page.Content = cursor.Current.Index(1).Value().String()
		d.Debugf("Get Page '%s'", title)
		return
	}

	d.Error(err)

	return

}

func (d *db) InitPage(title string, content string) error {
	_, err := d.page.InsertOne(context.TODO(), bson.D{{d.PageCollection.Title, title}, {d.PageCollection.Body, content}})
	d.Debugf("Initialize Page '%s'", title)
	return err
}

func (d *db) SetContent(title string, content string) error {
	_, err := d.page.UpdateByID(context.TODO(), title, bson.D{{"$set", bson.D{{d.PageCollection.Body, content}}}})
	d.Debugf("Amend Page '%s'", title)
	return err
}
