package main

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"modoowiki/db"
	"modoowiki/model/config"
)

type wiki struct {
	db DB
	config.Web
	*iris.Application
}

func Init(conf config.Config, level golog.Level) (w *wiki, err error) {

	w = new(wiki)

	w.db, err = db.New(conf.DB, level)
	if err != nil {
		return
	}

	w.Web = conf.Web
	w.Application = iris.New()
	w.Logger().Level = level

	return

}

func (w *wiki) Start() error {

	w.Get("/page/{key}", w.GetPage)
	w.Put("/page/{key}", w.PutPage)

	return w.Listen(w.Path())

}

func (w *wiki) Close() {
	w.db.Close()
	w.Logger().Info("Closing wiki web server...")
}

func (w *wiki) CommonLog(ctx context.Context) {
	w.Logger().Debugf("%s-> %s", ctx.RemoteAddr(), ctx.FullRequestURI())
}

func (w *wiki) GetPage(ctx context.Context) {

	w.CommonLog(ctx)
	key := ctx.Params().Get("key")

	page, err := w.db.GetPage(key)
	if err == nil {
		_, err = ctx.WriteString(page.String())
		w.Logger().Debug("Finish getting page")
	}

	if err != nil {
		w.Logger().Error(err)
	}

}

func (w *wiki) PutPage(ctx context.Context) {

	w.CommonLog(ctx)
	key := ctx.Params().Get("key")

	err := w.db.InitPage(key)
	if err != nil {
		w.Logger().Error(err)
	}

	w.Logger().Debug("Finish putting page")

}

func (w *wiki) InitPage(key string) (err error) {

	err = w.db.InitPage(key)
	w.Logger().Debug("Finish initializing page")

	return

}
