package gin

import (
	"github.com/freecloudio/server/plugin/graphql"

	"github.com/gin-gonic/gin"
)

func (r *Router) buildGraphQLRoutes() {
	grp := r.engine.Group("/api/graphql")

	grp.POST("", gin.WrapH(graphql.GetGraphQLHandler(r.cfg, r.managers)))
	grp.GET("", gin.WrapH(graphql.GetGraphQLPlaygroundHandler(r.cfg)))
}
