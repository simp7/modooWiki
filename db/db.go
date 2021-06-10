package db

import (
	"context"
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
}

func New(conf config.DB) (*db, error) {

	w := new(db)
	w.DB = conf

	w.end = make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(w.DB.Path()))
	if err == nil {
		w.Database = client.Database(w.DB.Name)
		w.page = w.Collection(w.PageCollection.Name)
	}

	go func() {
		for {
			select {
			case <-w.end:
				if err = client.Disconnect(ctx); err != nil {
					panic(err)
				}
			}
		}
	}()

	return w, err

}

func (w *db) Close() {
	close(w.end)
}

func (w *db) GetPage(title string) (page model.Page, err error) {

	cursor, err := w.page.Find(context.TODO(), bson.D{{w.PageCollection.Title, title}})

	if err == nil && cursor.Next(context.TODO()) {
		page.Key = cursor.Current.Index(0).Value().String()
		page.Content = cursor.Current.Index(1).Value().String()
	}

	return

}

func (w *db) InitPage(title string, content string) error {
	_, err := w.page.InsertOne(context.TODO(), bson.D{{w.PageCollection.Title, title}, {w.PageCollection.Body, content}})
	return err
}

func (w *db) SetContent(title string, content string) error {
	_, err := w.page.UpdateByID(context.TODO(), title, bson.D{{"$set", bson.D{{w.PageCollection.Body, content}}}})
	return err
}
