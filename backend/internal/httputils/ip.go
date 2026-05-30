package httputils

import (
	"context"
	"net"
)

type ctxKey string

const RealIPKey ctxKey = "real_ip"

func WithRealIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, RealIPKey, ip)
}

func GetRealIPWithPort(ctx context.Context) (string, string) {
	ipWithPort, _ := ctx.Value(RealIPKey).(string)
	ip, port, _ := net.SplitHostPort(ipWithPort)
	return ip, port
}

func GetRealIP(ctx context.Context) string {
	ipWithPort, _ := ctx.Value(RealIPKey).(string)
	ip, _, _ := net.SplitHostPort(ipWithPort)
	return ip
}

func GetRealPort(ctx context.Context) string {
	ipWithPort, _ := ctx.Value(RealIPKey).(string)
	_, port, _ := net.SplitHostPort(ipWithPort)
	return port
}
