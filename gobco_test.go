package gobco

import (
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	ReportCoverage()
	ReportProfile("instrument.out")
	os.Exit(retCode)
}

func TestCount(t *testing.T) {
	type args struct {
		cond  bool
		true  int
		false int
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{name: "true", args: args{cond: true, true: 0, false: 0}, want: args{cond: true, true: 1, false: 0}},
		{name: "false", args: args{cond: false, true: 0, false: 0}, want: args{cond: false, true: 0, false: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.args.cond, &tt.args.true, &tt.args.false);
				got != tt.want.cond || tt.args.true != tt.want.true || tt.args.false != tt.want.false {
				t.Errorf("Count() = %v, want %v", tt.args, tt.want)
			}
		})
	}
}

func Test_calculateCoverage(t *testing.T) {
	type result struct {
		cnt   int
		total int
	}
	tests := []struct {
		name string
		args *Cov
		want result
	}{
		{name: "true covered", args: &Cov{TCount: []int{1}, FCount: []int{0}}, want: result{cnt: 1, total: 2}},
		{name: "false covered", args: &Cov{TCount: []int{0}, FCount: []int{1}}, want: result{cnt: 1, total: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnt, total := calculateCoverage(tt.args)
			if cnt != tt.want.cnt {
				t.Errorf("calculateCoverage() cnt = %v, want %v", cnt, tt.want)
			}
			if total != tt.want.total {
				t.Errorf("calculateCoverage() total = %v, want %v", total, tt.want)
			}
		})
	}
}

func TestReportProfile(t *testing.T) {
	tests := []struct {
		name string
		args string
		want error
	}{
		{name: "wrong file", args: "/root", want: errors.New("")},
		{name: "wrong file", args: "./tmp", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReportProfile(tt.args)
			if tt.want != nil && err == nil {
				t.Errorf("ReportProfile() err = %v, want %v", err, tt.want)
			}
			if tt.want == nil && err != nil {
				t.Errorf("ReportProfile() err = %v, want %v", err, tt.want)
			}
		})
	}
}
