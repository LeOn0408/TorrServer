package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"server/log"
	"server/torr"
	"server/torr/state"
	"server/web/api/utils"
)

func torrentUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer form.RemoveAll()

	save := len(form.Value["save"]) > 0
	var retList []*state.TorrentStatus

	for name, file := range form.File {
		log.TLogln("add torrent file", name)

		torrFile, err := file[0].Open()
		if err != nil {
			log.TLogln("error upload torrent:", err)
			continue
		}
		defer torrFile.Close()

		spec, err := utils.ParseFile(torrFile)
		if err != nil {
			log.TLogln("error upload torrent:", err)
			continue
		}

		tor, err := torr.AddTorrent(spec, "", "")
		if err != nil {
			log.TLogln("error upload torrent:", err)
			continue
		}
		if save {
			torr.SaveTorrentToDB(tor)
		}
		retList = append(retList, tor.Status())
	}
	c.JSON(200, retList)
}