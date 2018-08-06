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
			alias, err := generateAlias(url)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(url, alias)
			err = store.AddRoute(url, alias)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println("cached route", url)
		}(url, store)

		c.Next()
	}
}
