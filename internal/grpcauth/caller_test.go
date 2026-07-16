package grpcauth

import "testing"

func TestCaller_IsStaff(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{"admin", []string{"user", "admin"}, true},
		{"waiter", []string{"waiter"}, true},
		{"restaurant admin", []string{"restaurantAdmin"}, true},
		{"plain user", []string{"user"}, false},
		{"no roles", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Caller{Roles: tt.roles}).IsStaff(); got != tt.want {
				t.Errorf("IsStaff(%v) = %v, want %v", tt.roles, got, tt.want)
			}
		})
	}
}

func TestCaller_OwnsOrStaff(t *testing.T) {
	tests := []struct {
		name   string
		caller Caller
		owner  string
		want   bool
	}{
		{"owner", Caller{ID: "u1"}, "u1", true},
		{"different user", Caller{ID: "u1"}, "u2", false},
		{"staff, not owner", Caller{ID: "u1", Roles: []string{"admin"}}, "u2", true},
		{"empty caller", Caller{}, "u1", false},
		{"empty caller and owner", Caller{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.caller.OwnsOrStaff(tt.owner); got != tt.want {
				t.Errorf("OwnsOrStaff(%q) = %v, want %v", tt.owner, got, tt.want)
			}
		})
	}
}
