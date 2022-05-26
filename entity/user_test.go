package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	got, err := NewUser("JS", "John", "Smith", "passwordHash$123654789", false)
	if err != nil {
		t.Errorf("NewUser() unexpected error = %v", err)
		return
	}

	if dif := time.Since(got.CreatedAt).Seconds(); dif < 0 || dif > 1 {
		t.Errorf("time.Since(NewUser().CreatedAt) = %v, want <1s", dif)
	} else {
		got.CreatedAt = time.Time{}
	}

	want := User{Entity{0, got.CreatedAt, got.UpdatedAt}, "JS", "John", "Smith", "passwordHash$123654789", false}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewUser() = %v, want %v", got, want)
	}

}

func TestNewUser_CreatedAt(t *testing.T) {
	got, err := NewUser("JS", "John", "Smith", "passwordHash$123654789", false)
	if err != nil {
		t.Errorf("NewUser() unexpected error = %v", err)
		return
	}

	if dif := time.Since(got.CreatedAt).Seconds(); dif < 0 || dif > 1 {
		t.Errorf("time.Since(NewUser().CreatedAt) = %v, want <1s", dif)
	}
}

func TestNewUser_Invalid(t *testing.T) {
	type args struct {
		username     string
		passwordHash string
		firstName    string
		lastName     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"missing username",
			args{"", "passwordHash$123654789", "John", "Smith"},
		},
		{
			"missing password hash",
			args{"JS", "", "John", "Smith"},
		},
		{
			"missing first name",
			args{"JS", "passwordHash$123654789", "", "Smith"},
		},
		{
			"missing last name",
			args{"JS", "passwordHash$123654789", "John", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.args.username, tt.args.firstName, tt.args.lastName, tt.args.passwordHash, false)
			if err == nil {
				t.Errorf("NewUser() expected a validation error but got nothing")
				return
			}
		})
	}
}
