package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"rdbviewer/defs"
	"rdbviewer/shared"
)

func (db *database) Search(ctx context.Context, req *shared.SearchReq) (*shared.SearchRes, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	out := &shared.SearchRes{}

	queryKeywords := getTextKeywords(req.Query, 1)
	if len(queryKeywords) == 0 {
		return out, nil
	}

	// search by TTH
	if len(queryKeywords) == 1 && len(firstKey(queryKeywords)) == 39 {
		tth := strings.ToUpper(firstKey(queryKeywords))

		if b := func() *defs.RadioBootleg {
			for _, b := range db.data.Bootlegs {
				for _, f := range b.Files {
					if f.TTH == tth {
						return b
					}
				}
			}
			return nil
		}(); b != nil {
			out.Items = append(out.Items, &shared.SearchResItem{
				Url:   "/bootlegs/" + b.Id,
				Title: fmt.Sprintf("%s", b.Name),
			})
		}

		// search by keywords
	} else {
		type Score struct {
			Show  *defs.RadioShow
			Score int
		}
		var scores []Score

		for id, kws := range db.showKeywords {
			score := 0
			for qw := range queryKeywords {
				score += func() int {
					for kw := range kws {
						if strings.Contains(kw, qw) {
							return 1
						}
					}
					return 0
				}()
			}
			if score >= len(queryKeywords) {
				scores = append(scores, Score{
					Show:  db.data.Shows[id],
					Score: score,
				})
			}
		}

		// sort by score, date
		sort.Slice(scores, func(i, j int) bool {
			if scores[i].Score != scores[j].Score {
				return scores[i].Score > scores[j].Score
			}
			return scores[i].Show.Date > scores[j].Show.Date
		})

		// keep first 10
		if len(scores) > 10 {
			scores = scores[:10]
		}

		for _, sco := range scores {
			sd, _ := time.Parse("2006-01-02", sco.Show.Date)
			out.Items = append(out.Items, &shared.SearchResItem{
				Url:   "/shows/" + sco.Show.Id,
				Title: fmt.Sprintf("%s, %s, %s, %s", shared.LabelArtist(sco.Show), sd.Format("2 January 2006"), sco.Show.City, strings.ToUpper(sco.Show.CountryCode)),
			})
		}
	}

	return out, nil
}
