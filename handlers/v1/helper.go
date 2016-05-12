package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func min(v1, v2 int64) int64 {
	if v1 < v2 {
		return v1
	}
	return v2
}

func parsePagination(c *gin.Context) (limit, offset int64, err error) {
	// parse limit
	queryLimit := c.DefaultQuery("limit", "10")
	limit, err = strconv.ParseInt(queryLimit, 10, 64)
	if err != nil {
		return
	}
	limit = min(limit, 100)

	// parse offset
	queryOffset := c.DefaultQuery("offset", "0")
	offset, err = strconv.ParseInt(queryOffset, 10, 64)
	if err != nil {
		return
	}

	return
}
