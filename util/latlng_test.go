package util

import (
	"testing"
)

func TestLatLng_DistanceTo(t *testing.T) {
	type args struct {
		newL LatLng
	}
	tests := []struct {
		name string
		l    LatLng
		args args
		want float64
	}{
		{
			name: "same latlng",
			l:    LatLng{24.182161901582095, 120.69032631143352},
			args: args{
				newL: LatLng{24.182161901582095, 120.69032631143352},
			},
			want: 0,
		},
		{
			name: "case 1",
			l:    LatLng{24.191234601023595, 120.6606317994254},
			args: args{
				newL: LatLng{24.191225426029785, 120.66066063317214},
			},
			want: 3.0974469364978376,
		},
		{
			name: "case 2",
			l:    LatLng{24.191234601023595, 120.6606317994254},
			args: args{
				newL: LatLng{24.191022546928178, 120.66049025856091},
			},
			want: 27.606064772337692,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.DistanceTo(tt.args.newL); got != tt.want {
				t.Errorf("LatLng.DistanceTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatLng_IsWithinRadiusOf(t *testing.T) {
	type args struct {
		latLng LatLng
		radius float64
	}
	tests := []struct {
		name string
		l    LatLng
		args args
		want bool
	}{
		{
			name: "same latlng",
			l:    LatLng{24.182161901582095, 120.69032631143352},
			args: args{
				latLng: LatLng{24.182161901582095, 120.69032631143352},
				radius: 5,
			},
			want: true,
		},
		{
			name: "in within 5m",
			l:    LatLng{24.19947356078559, 120.68417057318712},
			args: args{
				latLng: LatLng{24.19946316313104, 120.68415917379888},
				radius: 5,
			},
			want: true,
		},
		{
			name: "not in within 5m",
			l:    LatLng{24.19947356078559, 120.68417057318712},
			args: args{
				latLng: LatLng{24.19960449049318, 120.683934354949},
				radius: 5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IsWithinRadiusOf(tt.args.latLng, tt.args.radius); got != tt.want {
				t.Errorf("LatLng.IsWithinRadiusOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
