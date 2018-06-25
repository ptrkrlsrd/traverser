package acache

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/utilities/pkg/ucrypt"
	"github.com/ptrkrlsrd/utilities/pkg/unet"
)

func generateAlias(url string) string {
	splitUrl, _ := unet.SplitUrl(url)

	if len(splitUrl) > 1 {
		return splitUrl[1]
	}

	hash := ucrypt.MD5Hash(url)
	return hash
}

func CacheRoute(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = c.Request
		var url = req.URL.String()

		log.Println("cached route", url)

		go func(url string) {
			if hasRoute, err := store.HasRoute(url); err == nil && !hasRoute {
				var alias = generateAlias(url)
				err = store.AddRoute(url, alias)
				if err != nil {
					log.Println(err)
				}
			}
		}(url)

		c.Next()
	}
}
