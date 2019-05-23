package phonenumber

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		args    args
		want    PhoneNumber
		wantErr bool
	}{
		{name: "Local", args: args{"754-3010"}, want: PhoneNumber{Local: [4]byte{7, 54, 30, 10}}},
		{name: "Domestic", args: args{"(541) 754-3010"}, want: PhoneNumber{OperatorCode: 541, Local: [4]byte{7, 54, 30, 10}}},
		{name: "International", args: args{"+1-541-754-3010"}, want: PhoneNumber{CountyCode: 1, OperatorCode: 541, Local: [4]byte{7, 54, 30, 10}}},
		{name: "Dialed in the US", args: args{"1-541-754-3010"}, want: PhoneNumber{CountyCode: 1, OperatorCode: 541, Local: [4]byte{7, 54, 30, 10}}},
		{name: "Invalid input", args: args{"fooNumber"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
