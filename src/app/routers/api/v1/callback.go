package v1

import (
	"github.com/devops-salt/src/log"
	"github.com/gin-gonic/gin"
)

func Callback(c *gin.Context)  {
	d, err := c.GetRawData()
	if err!=nil {
		log.Error("%s", err)
	}
	log.Error("%s", string(d))
	c.JSON(200, "")
}
