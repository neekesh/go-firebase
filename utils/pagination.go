package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int
	Cursor string
}

func BuildPagination(c *gin.Context) Pagination {
	limitStr := c.Query("limit")
	cursor := c.Query("cursor")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 1
	}

	return Pagination{
		Limit:  limit,
		Cursor: cursor,
	}

}
