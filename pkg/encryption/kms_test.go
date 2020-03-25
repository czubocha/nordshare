package encryption

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/service/kms"
	"testing"
)

const encrypted = "encrypted"
const decrypted = "decrypted"

func Test_service_Encrypt(t *testing.T) {
	type fields struct {
		client crypter
	}
	type args struct {
		content *[]byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "successful encryption",
			fields: fields{
				client: crypterMock{},
			},
			args: args{
				content: &[]byte{1},
			},
			wantContent: []byte(encrypted),
			wantErr:     false,
		},
		{
			name: "unsuccessful encryption",
			fields: fields{
				client: crypterMock{errors.New("encryption error")},
			},
			args: args{
				content: &[]byte{1},
			},
			wantContent: []byte{1},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				client: tt.fields.client,
			}
			if err := s.Encrypt(tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(*tt.args.content, tt.wantContent) {
				t.Errorf("Encrypt() content = %v, want %v", *tt.args.content, tt.wantContent)
			}
		})
	}
}

func Test_service_Decrypt(t *testing.T) {
	type fields struct {
		client crypter
	}
	type args struct {
		content *[]byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "successful decryption",
			fields: fields{
				client: crypterMock{},
			},
			args: args{
				content: &[]byte{1},
			},
			wantContent: []byte(decrypted),
			wantErr:     false,
		},
		{
			name: "unsuccessful decryption",
			fields: fields{
				client: crypterMock{errors.New("decryption error")},
			},
			args: args{
				content: &[]byte{1},
			},
			wantContent: []byte{1},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService("", tt.fields.client)
			if err := s.Decrypt(tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(*tt.args.content, tt.wantContent) {
				t.Errorf("Decrypt() content = %v, want %v", *tt.args.content, tt.wantContent)
			}
		})
	}
}

type crypterMock struct {
	error
}

func (cm crypterMock) Encrypt(*kms.EncryptInput) (*kms.EncryptOutput, error) {
	if cm.error != nil {
		return &kms.EncryptOutput{}, cm.error
	}
	return &kms.EncryptOutput{CiphertextBlob: []byte(encrypted)}, cm.error
}
func (cm crypterMock) Decrypt(*kms.DecryptInput) (*kms.DecryptOutput, error) {
	return &kms.DecryptOutput{Plaintext: []byte(decrypted)}, cm.error
}
