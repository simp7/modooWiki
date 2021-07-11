package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/logger"
	"io"
	"modoowiki/db"
	"modoowiki/model/config"
	"os"
)

type wiki struct {
	db DB
	config.Web
	*logger.Logger
	*gin.Engine
}

func Init(conf config.Config, lg *logger.Logger) (w *wiki, err error) {

	w = new(wiki)

	w.db, err = db.New(conf.DB, lg)
	if err != nil {
		return
	}

	w.Web = conf.Web
	w.Engine = gin.New()
	w.Logger = lg

	return

}

func (w *wiki) Start() error {

	w.GET("/page/:key", w.GetPage)
	w.PUT("/page/:key", w.PutPage)

	return w.Run(w.Path())

}

func (w *wiki) Close() {
	w.db.Close()
	w.Info("Closing wiki web server...")
}

func (w *wiki) CommonLog(ctx *gin.Context) {
	ip, ok := ctx.RemoteIP()
	if ok {
		w.Info("%s-> %s", ip, ctx.Request.URL)
	}
}

func (w *wiki) GetPage(ctx *gin.Context) {

	w.CommonLog(ctx)
	key, ok := ctx.Params.Get("key")

	if ok {
		page, err := w.db.GetPage(key)

		if err == nil {
			_, err = ctx.Writer.WriteString(page.String())
			w.Info("Finish getting page")
		}

		if err != nil {
			w.Error(err)
		}
		return

	}

	if indexPage, err := w.IndexPage(); err == nil {
		_, err = ctx.Writer.Write(indexPage)
		if err != nil {
			w.Error(err)
		}
	}

	return

}

func (w *wiki) PutPage(ctx *gin.Context) {

	w.CommonLog(ctx)
	key, ok := ctx.Params.Get("key")

	if ok {

		err := w.db.InitPage(key)
		if err != nil {
			w.Error(err)
		}

		w.Info("Finish putting page")

	}

}

func (w *wiki) InitPage(key string) (err error) {

	err = w.db.InitPage(key)
	w.Info("Finish initializing page")

	return

}

func (w *wiki) IndexPage() ([]byte, error) {

	var err error
	var file io.Reader

	if file, err = os.Open("indexPage.html"); err == nil {
		return io.ReadAll(file)
	}

	return nil, err

}
