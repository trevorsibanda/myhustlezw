package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

func apiError(ctx *gin.Context, reason interface{}) {
	util.ApiError(ctx, fmt.Sprintf("%v", reason))
}
