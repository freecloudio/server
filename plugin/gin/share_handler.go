package gin

import (
	"net/http"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/gin-gonic/gin"
)

func (r *Router) buildShareRoutes() {
	grp := r.engine.Group("/api/share")

	grp.POST("", r.createShare)
}

func (r *Router) createShare(c *gin.Context) {
	authContext := getAuthContext(c)

	share := &models.Share{}
	err := c.BindJSON(share)
	if err != nil {
		fcerr := fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	created, fcerr := r.managers.Share.CreateShare(authContext, share)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"created": created})
}
