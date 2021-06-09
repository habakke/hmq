package broker

import (
	"github.com/gin-gonic/gin"
	"github.com/habakke/hmq/metrics"
)

func InitHTTP(b *Broker) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	b.metrics = &metrics.Manager{}
	b.metrics.Init(router)
	_ = b.metrics.Add(metrics.MetricNumberOfMessages, "Total number of packets received")
	_ = b.metrics.Add(metrics.MetricNumberOfClients, "Total number of clients connected")

	router.DELETE("api/v1/connections/:clientid", func(c *gin.Context) {
		clientid := c.Param("clientid")
		cli, ok := b.clients.Load(clientid)
		if ok {
			conn, succss := cli.(*client)
			if succss {
				conn.Close()
			}
		}
		resp := map[string]int{
			"code": 0,
		}
		c.JSON(200, &resp)
	})

	_ = router.Run(":" + b.config.HTTPPort)
}
