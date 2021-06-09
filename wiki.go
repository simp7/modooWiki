package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"simpleWiki/db"
	"simpleWiki/model/config"
)

type wiki struct {
	db DB
	config.Web
	*iris.Application
}

func Init(conf config.Config) (w *wiki, err error) {

	w = new(wiki)

	w.db, err = db.New(conf.DB)
	if err != nil {
		return
	}

	w.Web = conf.Web
	w.Application = iris.New()

	return

}

func (w *wiki) Start() error {
	w.Get("/page/{key}", w.GetPage)
	w.Post("/page/{key}", w.PostPage)
	return w.Listen(w.Path())
}

func (w *wiki) End() {
	w.db.Close()
}

func (w *wiki) CommonLog(ctx context.Context) {
	w.Logger().Infof("%s", ctx.RemoteAddr())
	w.Logger().Debugf("%s", ctx.FullRequestURI())
}

func (w *wiki) GetPage(ctx context.Context) {

	w.CommonLog(ctx)
	key := ctx.Params().Get("key")

	page, err := w.db.GetPage(key)
	if err == nil {
		_, err = ctx.WriteString(page.Content)
	}

	if err != nil {
		w.Logger().Error(err)
	}

}

func (w *wiki) PostPage(ctx context.Context) {

	w.CommonLog(ctx)
	key := ctx.Params().Get("key")
	content := ctx.PostValue("content")

	err := w.db.InitPage(key, content)
	if err != nil {
		w.Logger().Error(err)
	}

}

func (w *wiki) InitPage(key string, content string) (err error) {
	err = w.db.InitPage(key, content)
	return
}

func (w *wiki) SetContent(key string, content string) {
	w.db.SetContent(key, content)
}