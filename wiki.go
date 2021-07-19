package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/logger"
	"modoowiki/db"
	"modoowiki/model/config"
	"net/http"
)

//go:embed templates
var asset embed.FS

type wiki struct {
	db DB
	config.Web
	*logger.Logger
	*gin.Engine
}

func Init(conf config.Config, lg *logger.Logger) (w *wiki, err error) {

	w = new(wiki)

	w.db, err = db.Mongo(conf.DB, lg)
	if err != nil {
		return
	}

	w.Web = conf.Web
	w.Engine = gin.New()
	w.Logger = lg

	return

}

func (w *wiki) Start() error {

	w.GET("/", w.Index)
	w.GET("/page/:key", w.GetPage)
	w.PUT("/page/:key", w.PutPage)
	w.DELETE("/page/:key", w.DeletePage)

	w.LoadHTMLGlob("templates/html/*")

	return w.Run(w.Path())

}

func (w *wiki) Close() {
	w.db.Close()
	w.Info("Closing wiki web server...")
}

func (w *wiki) CommonLog(ctx *gin.Context) {
	ip, ok := ctx.RemoteIP()
	if ok {
		w.Infof("%s-> %s", ip, ctx.Request.URL)
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

func (w *wiki) DeletePage(ctx *gin.Context) {
	w.CommonLog(ctx)
	key, ok := ctx.Params.Get("key")

	if ok {
		err := w.db.DeletePage(key)
		if err != nil {
			w.Error(err)
		}
	}

}

func (w *wiki) InitPage(key string) (err error) {

	err = w.db.InitPage(key)
	w.Info("Finish initializing page")

	return

}

func (w *wiki) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "Home"})
}
