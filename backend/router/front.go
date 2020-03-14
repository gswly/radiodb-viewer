package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onPageFront(c *gin.Context) {
	front, err := h.dbClient.Front(context.Background(), &shared.FrontReq{})
	if err != nil {
		panic(err)
		http.Error(c.Writer, "500 internal server error", http.StatusInternalServerError)
		return
	}

	g, _ := time.Parse(time.RFC3339, front.Stats.Generated)

	limitSize := func(in string) string {
		if len(in) > 55 {
			in = in[:55] + "..."
		}
		return in
	}

	ginTemplate(c, h.frameWrapper(c, FrameConf{
		Title: "",
		Class: "front",
		Content: templateRender(h.templates["front"], gin.H{
			"ShowsLast": func() (ret []gin.H) {
				for _, s := range front.LastShows {
					da, _ := time.Parse("2006-01-02", s.Date)
					ret = append(ret, gin.H{
						"Id":      s.Id,
						"Title":   fmt.Sprintf("%s, %s", shared.LabelArtist(s), da.Format("2 January 2006")),
						"Country": s.City + ", " + strings.ToUpper(s.CountryCode),
						"Tour":    s.Tour,
					})
				}
				return ret
			}(),
			"BootlegsLast": func() (ret []gin.H) {
				for _, b := range front.LastBootlegs {
					s := front.LastBootlegShows[b.Show]
					ret = append(ret, gin.H{
						"Id":   b.Id,
						"Name": limitSize(b.Name),
						"Tour": s.Tour,
					})
				}
				return ret
			}(),
			"Generated":    g.Format("2 January 2006"),
			"BootlegCount": front.Stats.BootlegCount,
			"ShowCount":    front.Stats.ShowCount,
			"ShareSize":    humanize.Bytes(front.Stats.ShareSize),
		}),
	}))
}