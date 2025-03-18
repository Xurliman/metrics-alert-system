package middlewares

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"net"
)

type NetMiddleware struct {
	cidrStr string
}

func (m NetMiddleware) Handle(ctx *gin.Context) {
	if m.cidrStr == "" {
		ctx.Next()
		return
	}

	ipStr := ctx.Request.Header.Get("X-Real-IP")
	if !isIPInSubnet(ipStr, m.cidrStr) {
		utils.JSONForbiddenError(ctx, fmt.Errorf("access denied %s", ipStr))
		return
	}
	ctx.Next()
}

func NewNetMiddleware(cidrStr string) Middleware {
	return &NetMiddleware{
		cidrStr: cidrStr,
	}
}

func isIPInSubnet(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	_, subnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}
	return subnet.Contains(ip)
}
