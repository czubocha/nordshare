package hash

import (
	"github/czubocha/nordshare"
	"testing"
)

func TestPasswords(t *testing.T) {
	readPassword, writePassword, otherPassword := []byte("abc"), []byte("def"), []byte("incorrect")
	type args struct {
		passwordToRead, passwordToWrite           []byte
		shouldHasReadAccess, shouldHasWriteAccess bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "read pass for read ok, write pass for write ok",
			args: args{
				passwordToRead:       readPassword,
				passwordToWrite:      writePassword,
				shouldHasReadAccess:  true,
				shouldHasWriteAccess: true,
			},
		},
		{
			name: "write pass for read ok, write pass for write ok",
			args: args{
				passwordToRead:       writePassword,
				passwordToWrite:      writePassword,
				shouldHasReadAccess:  true,
				shouldHasWriteAccess: true,
			},
		},
		{
			name: "read pass for read ok, read pass for write ko",
			args: args{
				passwordToRead:       readPassword,
				passwordToWrite:      readPassword,
				shouldHasReadAccess:  true,
				shouldHasWriteAccess: false,
			},
		},
		{
			name: "other pass for read ko, write pass for write ok",
			args: args{
				passwordToRead:       otherPassword,
				passwordToWrite:      writePassword,
				shouldHasReadAccess:  false,
				shouldHasWriteAccess: true,
			},
		},
		{
			name: "read pass for read ok, other pass for write ko",
			args: args{
				passwordToRead:       readPassword,
				passwordToWrite:      otherPassword,
				shouldHasReadAccess:  true,
				shouldHasWriteAccess: false,
			},
		},
		{
			name: "other pass for read ko, other pass for write ko",
			args: args{
				passwordToRead:       otherPassword,
				passwordToWrite:      otherPassword,
				shouldHasReadAccess:  false,
				shouldHasWriteAccess: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nordshare.Note{
				ReadPassword:  readPassword,
				WritePassword: writePassword,
			}
			if err := Passwords(n); err != nil {
				t.Errorf("Passwords() error = %v", err)
			}
			if HasReadAccess(*n, tt.args.passwordToRead) != tt.args.shouldHasReadAccess {
				t.Errorf("Passwords() should have read access")
			}
			if HasWriteAccess(*n, tt.args.passwordToWrite) != tt.args.shouldHasWriteAccess {
				t.Errorf("Passwords() should have write access")
			}
		})
	}
}
