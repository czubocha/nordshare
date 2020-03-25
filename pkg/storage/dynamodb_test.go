package storage

import (
	"testing"
	"time"
)

func Test_isExpired(t *testing.T) {
	type args struct {
		ttl int64
		now time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "expired",
			args: args{
				ttl: 1585091931, // 03/24/2020 @ 11:18pm (UTC)
				now: time.Date(2020, 3, 25, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "not expired",
			args: args{
				ttl: 1585091931, // 03/24/2020 @ 11:18pm (UTC)
				now: time.Date(2020, 3, 24, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExpired(tt.args.ttl, tt.args.now); got != tt.want {
				t.Errorf("isExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateTTL(t *testing.T) {
	type args struct {
		minutesToExpire int64
		now             time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "default for tomorrow (24 hours) when 0 or less",
			args: args{
				minutesToExpire: 0,
				now:             time.Date(2020, 3, 24, 0, 0, 0, 0, time.UTC),
			},
			want: 1585094400, // 03/25/2020 @ 12:00am (UTC)
		},
		{
			name: "when positive minutes provided then minutes added to now",
			args: args{
				minutesToExpire: 100,
				now:             time.Date(2020, 3, 24, 0, 0, 0, 0, time.UTC),
			},
			want: 1585014000, // 03/24/2020 @ 1:40am (UTC)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTTL(tt.args.minutesToExpire, tt.args.now); got != tt.want {
				t.Errorf("calculateTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRemainingMinutes(t *testing.T) {
	type args struct {
		unixTime int64
		now      time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "5 minutes remain to unixTime",
			args: args{
				unixTime: 1585014000, // 03/24/2020 @ 1:40am (UTC)
				now:      time.Date(2020, 3, 24, 1, 35, 0, 0, time.UTC),
			},
			want: 5,
		},
		{
			name: "0 when now is later than unixTime",
			args: args{
				unixTime: 1585014000, // 03/24/2020 @ 1:40am (UTC)
				now:      time.Date(2020, 3, 25, 1, 50, 0, 0, time.UTC),
			},
			want: 0,
		},
		{
			name: "0 when the same time",
			args: args{
				unixTime: 1585014000, // 03/24/2020 @ 1:40am (UTC)
				now:      time.Date(2020, 3, 24, 1, 40, 0, 0, time.UTC),
			},
			want: 0,
		},
		{
			name: "1 when remain less than minute",
			args: args{
				unixTime: 1585014000, // 03/24/2020 @ 1:40am (UTC)
				now:      time.Date(2020, 3, 24, 1, 39, 1, 0, time.UTC),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRemainingMinutes(tt.args.unixTime, tt.args.now); got != tt.want {
				t.Errorf("getRemainingMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
