package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func min(v1, v2 int64) int64 {
	if v1 < v2 {
		return v1
	}
	return v2
}

func parsePagination(c *gin.Context) (isSince bool, separator, limit int64, err error) {
	// parse limit
	queryLimit := c.DefaultQuery("limit", "10")
	limit, err = strconv.ParseInt(queryLimit, 10, 64)
	if err != nil {
		return
	}
	limit = min(limit, 100)

	// parse since
	querySince := c.Query("since")
	queryUntil := c.Query("until")
	switch {
	case querySince != "":
		isSince = true
		separator, err = strconv.ParseInt(querySince, 10, 64)
	case queryUntil != "":
		isSince = false
		separator, err = strconv.ParseInt(queryUntil, 10, 64)
	default:
		err = errors.New("since or until should present")
	}

	return
}
