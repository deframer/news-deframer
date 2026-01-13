package syncer

import (
	"fmt"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
)

// applyFancyRatingText call this after you have the Title = TitleCorrected
func applyFancyRatingText(item *gofeed.Item, res *database.ThinkResult, thinkRating float64, language string) {
	if item == nil || res == nil {
		return
	}

	item.Title = fmt.Sprintf("%s %s", scoreToStars(thinkRating), item.Title)
}
