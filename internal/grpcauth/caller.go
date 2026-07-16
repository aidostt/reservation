// Package grpcauth propagates the authenticated caller's identity from the
// gateway (over gRPC metadata) into the request context, so the reservation
// domain can enforce ownership itself rather than trusting the caller.
package grpcauth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	mdCallerID    = "caller-id"
	mdCallerRoles = "caller-roles"
)

// Staff roles may act on any reservation regardless of ownership.
const (
	RoleAdmin           = "admin"
	RoleRestaurantAdmin = "restaurantAdmin"
	RoleWaiter          = "waiter"
)

type Caller struct {
	ID    string
	Roles []string
}

func (c Caller) IsStaff() bool {
	for _, r := range c.Roles {
		switch r {
		case RoleAdmin, RoleRestaurantAdmin, RoleWaiter:
			return true
		}
	}
	return false
}

// OwnsOrStaff reports whether the caller owns the resource identified by
// ownerID, or is a staff member.
func (c Caller) OwnsOrStaff(ownerID string) bool {
	return c.ID != "" && c.ID == ownerID || c.IsStaff()
}

type callerKey struct{}

// UnaryInterceptor reads the caller identity from gRPC metadata into the
// context. A missing identity yields the zero Caller, which owns nothing and is
// not staff, so handlers deny by default.
func UnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var caller Caller
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v := md.Get(mdCallerID); len(v) > 0 {
			caller.ID = v[0]
		}
		if v := md.Get(mdCallerRoles); len(v) > 0 && v[0] != "" {
			caller.Roles = strings.Split(v[0], ",")
		}
	}
	return handler(context.WithValue(ctx, callerKey{}, caller), req)
}

// FromContext returns the caller extracted by UnaryInterceptor.
func FromContext(ctx context.Context) Caller {
	c, _ := ctx.Value(callerKey{}).(Caller)
	return c
}
