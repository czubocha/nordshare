package api

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	type args struct {
		status int
		body   []interface{}
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "200 without body",
			args: args{
				status: http.StatusOK,
				body:   nil,
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "",
			},
			wantErr: false,
		},
		{
			name: "200 with struct body",
			args: args{
				status: http.StatusOK,
				body: []interface{}{struct {
					Exported     string
					WithTag      bool `json:"withTag"`
					ZeroValue    byte
					ZeroValueMap map[string]string
					ByteArray    []byte
					unexported   int
				}{
					Exported:   "test",
					WithTag:    true,
					ByteArray:  []byte("a"),
					unexported: 123,
				}},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"Exported":"test","withTag":true,"ZeroValue":0,"ZeroValueMap":null,"ByteArray":"YQ=="}`,
			},
			wantErr: false,
		},
		{
			name: "404 with just string body",
			args: args{
				status: http.StatusNotFound,
				body:   []interface{}{"error"},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       `"error"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := NewResponse(tt.args.status, tt.args.body...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResponse() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("NewResponse() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_getHeaderIncasesensible(t *testing.T) {
	type args struct {
		headers map[string]string
		key     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "exact match",
			args: args{
				headers: map[string]string{"Password": "abc"},
				key:     "Password",
			},
			want: "abc",
		},
		{
			name: "uppercase",
			args: args{
				headers: map[string]string{"PASSWORD": "abc"},
				key:     "Password",
			},
			want: "abc",
		},
		{
			name: "lowercase",
			args: args{
				headers: map[string]string{"password": "abc"},
				key:     "Password",
			},
			want: "abc",
		},
		{
			name: "nonexistent key",
			args: args{
				headers: map[string]string{"test": "abc"},
				key:     "Password",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHeaderIncasesensible(tt.args.headers, tt.args.key); got != tt.want {
				t.Errorf("GetHeaderIncasesensible() = %v, want %v", got, tt.want)
			}
		})
	}
}
