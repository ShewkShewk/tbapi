package tb_api

import (
	"errors"
	"reflect"
	"testing"
)

func Test_newBuilder(t *testing.T) {
	tests := []struct {
		name string
		want Builder
	}{
		{"NewBuilder returns an empty builder", Builder{
			ApiUrl:   "",
			Username: "",
			Password: "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBuilder(); !reflect.DeepEqual(*got, tt.want) { // Dereference got with *got
				t.Errorf("NewBuilder() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestBuilder_withApiUrl(t *testing.T) {
	type fields struct {
		ApiUrl   string
		Username string
		Password string
	}
	type args struct {
		apiUrl string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "WithApiUrl returns a builder with the URL set",
			fields: fields{
				ApiUrl:   "",
				Username: "",
				Password: "",
			},
			args: args{
				apiUrl: "https://www.tabroom.com",
			},
			want: &Builder{
				ApiUrl:   "https://www.tabroom.com",
				Username: "",
				Password: "",
			},
		},
		{
			name: "WithApiUrl will overwrite existing URL",
			fields: fields{
				ApiUrl:   "https://www.tabroom.com",
				Username: "",
				Password: "",
			},
			args: args{
				apiUrl: "https://www.other.tabroom.com",
			},
			want: &Builder{
				ApiUrl:   "https://www.other.tabroom.com",
				Username: "",
				Password: "",
			},
		},
		{
			name: "WithApiUrl will not effect other fields",
			fields: fields{
				ApiUrl:   "",
				Username: "a_username_value",
				Password: "a_password_value",
			},
			args: args{
				apiUrl: "https://www.tabroom.com",
			},
			want: &Builder{
				ApiUrl:   "https://www.tabroom.com",
				Username: "a_username_value",
				Password: "a_password_value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				ApiUrl:   tt.fields.ApiUrl,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
			if got := b.WithApiUrl(tt.args.apiUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithApiUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_withUsername(t *testing.T) {
	type fields struct {
		ApiUrl   string
		Username string
		Password string
	}
	type args struct {
		username string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "WithUsername will populate the username field",
			fields: fields{
				ApiUrl:   "",
				Username: "",
				Password: "",
			},
			args: args{
				username: "example_username",
			},
			want: &Builder{
				ApiUrl:   "",
				Username: "example_username",
				Password: "",
			},
		},
		{
			name: "WithUsername will overwrite existing username",
			fields: fields{
				ApiUrl:   "",
				Username: "example_username",
				Password: "",
			},
			args: args{
				username: "new_username",
			},
			want: &Builder{
				ApiUrl:   "",
				Username: "new_username",
				Password: "",
			},
		},
		{
			name: "WithUsername will not modify other fields",
			fields: fields{
				ApiUrl:   "https://tabroom.com",
				Username: "",
				Password: "example_password",
			},
			args: args{
				username: "example_username",
			},
			want: &Builder{
				ApiUrl:   "https://tabroom.com",
				Username: "example_username",
				Password: "example_password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				ApiUrl:   tt.fields.ApiUrl,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
			if got := b.WithUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithPassword(t *testing.T) {
	type fields struct {
		ApiUrl   string
		Username string
		Password string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name: "withPassword will populate the password field",
			fields: fields{
				ApiUrl:   "",
				Username: "",
				Password: "",
			},
			args: args{
				password: "example_password",
			},
			want: &Builder{
				ApiUrl:   "",
				Username: "",
				Password: "example_password",
			},
		},
		{
			name: "withPassword will overwrite the password field",
			fields: fields{
				ApiUrl:   "",
				Username: "",
				Password: "example_password",
			},
			args: args{
				password: "other_password",
			},
			want: &Builder{
				ApiUrl:   "",
				Username: "",
				Password: "other_password",
			},
		},
		{
			name: "withPassword will not modify other fields",
			fields: fields{
				ApiUrl:   "https://tabroom.com",
				Username: "example_username",
				Password: "",
			},
			args: args{
				password: "example_password",
			},
			want: &Builder{
				ApiUrl:   "https://tabroom.com",
				Username: "example_username",
				Password: "example_password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				ApiUrl:   tt.fields.ApiUrl,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
			if got := b.WithPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Build(t *testing.T) {
	type fields struct {
		ApiUrl   string
		Username string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
		wantApi *TabroomApi
	}{
		{
			name: "Builder with all fields populated will return a TabroomApi",
			fields: fields{
				ApiUrl:   "https://tabroom.com",
				Username: "example_username",
				Password: "example_password",
			},
			wantErr: nil,
			wantApi: &TabroomApi{},
		},
		{
			name: "Builder with url missing will return an error",
			fields: fields{
				ApiUrl:   "",
				Username: "example_username",
				Password: "example_password",
			},
			wantErr: errors.New("missing API URL in builder"),
			wantApi: nil,
		},
		{
			name: "Builder with username missing will return an error",
			fields: fields{
				ApiUrl:   "https://tabroom.com",
				Username: "",
				Password: "example_password",
			},
			wantErr: errors.New("missing username in builder"),
			wantApi: nil,
		},
		{
			name: "Builder with username missing will return an error",
			fields: fields{
				ApiUrl:   "https://tabroom.com",
				Username: "example_username",
				Password: "",
			},
			wantErr: errors.New("missing password in builder"),
			wantApi: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				ApiUrl:   tt.fields.ApiUrl,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
			got, got1 := b.Build()
			if !reflect.DeepEqual(got, tt.wantErr) {
				t.Errorf("Build() got = %v, want %v", got, tt.wantErr)
			}
			if !reflect.DeepEqual(got1, tt.wantApi) {
				t.Errorf("Build() got1 = %v, want %v", got1, tt.wantApi)
			}
		})
	}
}
