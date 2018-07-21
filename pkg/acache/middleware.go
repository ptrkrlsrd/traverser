package acache

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CacheRoute(store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = c.Request
		var url = req.RequestURI

		go func(url string, store *Store) {
			var alias = generateAlias(url)
			log.Println(url, alias)
			err := store.AddRoute(url, alias)
			if err != nil {
				log.Println(err)
			}

			log.Println("cached route", url)
		}(url, store)

		c.Next()
	}
}
