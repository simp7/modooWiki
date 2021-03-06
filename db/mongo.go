package db

import (
	"context"
	"github.com/google/logger"
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
	*logger.Logger
}

func Mongo(conf config.DB, lg *logger.Logger) (*db, error) {

	d := new(db)
	d.DB = conf

	d.Logger = lg

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

	cursor, err := d.page.Find(context.TODO(), bson.D{{d.PageCollection.Key, title}})

	if err == nil && cursor.Next(context.TODO()) {
		page.Root.Title = cursor.Current.Index(0).Value().String()
		page.Root.Content = cursor.Current.Index(1).Value().String()
		d.Infof("Get Page '%s'", title)
		return
	}

	d.Error("Page not found.")

	return

}

func (d *db) InitPage(title string) error {
	page := model.NewPage(title)
	_, err := d.page.InsertOne(context.TODO(), bson.D{{d.PageCollection.Key, title}, {d.PageCollection.Page, page}})
	d.Infof("Initialize Page '%s'", title)
	return err
}

func (d *db) SetContent(page model.Page) error {
	_, err := d.page.UpdateByID(context.TODO(), page.Key(), bson.D{{"$set", bson.D{{d.PageCollection.Page, page}}}})
	d.Infof("Amend Page '%s'", page.Key())
	return err
}

func (d *db) DeletePage(title string) error {
	_, err := d.page.DeleteOne(context.TODO(), bson.D{{d.PageCollection.Key, title}})
	return err
}
