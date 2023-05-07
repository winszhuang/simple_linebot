package util

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Paginate(t *testing.T) {
	// type string
	type args struct {
		data      []string
		pageIndex int
		pageSize  int
	}

	genStringDataList := func() []string {
		arr := make([]string, 50)
		for i := 0; i < 50; i++ {
			arr[i] = fmt.Sprintf("test%d", i+1)
		}
		return arr
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "String : Base",
			args: args{
				data:      genStringDataList(),
				pageIndex: 2,
				pageSize:  5,
			},
			want: []string{
				"test6",
				"test7",
				"test8",
				"test9",
				"test10",
			},
		},
		{
			name: "String : Exceed the limit",
			args: args{
				data:      genStringDataList(),
				pageIndex: 5,
				pageSize:  12,
			},
			want: []string{
				"test49",
				"test50",
			},
		},
		{
			name: "String : Exceed the limit2",
			args: args{
				data:      genStringDataList(),
				pageIndex: 6,
				pageSize:  15,
			},
			want: []string{},
		},
		{
			name: "String : zero pageIndex",
			args: args{
				data:      genStringDataList(),
				pageIndex: 0,
				pageSize:  15,
			},
			want: []string{},
		},
		{
			name: "String : negative pageIndex",
			args: args{
				data:      genStringDataList(),
				pageIndex: -5,
				pageSize:  15,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Paginate(tt.args.data, tt.args.pageIndex, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paginate() = %v, want %v", got, tt.want)
			}
		})
	}

	// type 2
	type args2 struct {
		data      []LatLng
		pageIndex int
		pageSize  int
	}

	genLatLngDataList := func() []LatLng {
		arr := make([]LatLng, 50)
		for i := 0; i < 50; i++ {
			arr[i] = LatLng{Lat: float64(i + 1), Lng: float64(i + 2)}
		}
		return arr
	}

	tests2 := []struct {
		name string
		args args2
		want []LatLng
	}{
		{
			name: "LatLng : Base",
			args: args2{
				data:      genLatLngDataList(),
				pageIndex: 2,
				pageSize:  5,
			},
			want: []LatLng{
				{Lat: 6, Lng: 7},
				{Lat: 7, Lng: 8},
				{Lat: 8, Lng: 9},
				{Lat: 9, Lng: 10},
				{Lat: 10, Lng: 11},
			},
		},
		{
			name: "LatLng : Exceed the limit",
			args: args2{
				data:      genLatLngDataList(),
				pageIndex: 5,
				pageSize:  12,
			},
			want: []LatLng{
				{Lat: 49, Lng: 50},
				{Lat: 50, Lng: 51},
			},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := Paginate(tt.args.data, tt.args.pageIndex, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paginate() = %v, want %v", got, tt.want)
			}
		})
	}
}
