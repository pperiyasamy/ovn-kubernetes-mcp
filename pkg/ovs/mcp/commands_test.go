package mcp

import (
	"testing"
)

func TestFilterLines(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		pattern string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty pattern returns all lines",
			lines:   []string{"line1", "line2", "line3"},
			pattern: "",
			want:    []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "simple string match",
			lines:   []string{"foo", "bar", "foo bar", "baz"},
			pattern: "foo",
			want:    []string{"foo", "foo bar"},
			wantErr: false,
		},
		{
			name:    "regex pattern match",
			lines:   []string{"table=0", "table=10", "priority=100", "table=5"},
			pattern: `table=\d+`,
			want:    []string{"table=0", "table=10", "table=5"},
			wantErr: false,
		},
		{
			name:    "no matches returns empty slice",
			lines:   []string{"line1", "line2", "line3"},
			pattern: "nomatch",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid regex pattern returns error",
			lines:   []string{"line1", "line2"},
			pattern: "[invalid",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty lines with pattern",
			lines:   []string{},
			pattern: "test",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "complex regex with multiple groups",
			lines:   []string{"cookie=0x0, table=0", "cookie=0x1, table=10", "priority=100"},
			pattern: `cookie=0x[0-9a-f]+`,
			want:    []string{"cookie=0x0, table=0", "cookie=0x1, table=10"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterLines(tt.lines, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("filterLines() got %d lines, want %d lines", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("filterLines()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestLimitLines(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		maxLines int
		want     []string
	}{
		{
			name:     "zero maxLines returns all lines",
			lines:    []string{"line1", "line2", "line3"},
			maxLines: 0,
			want:     []string{"line1", "line2", "line3"},
		},
		{
			name:     "negative maxLines returns all lines",
			lines:    []string{"line1", "line2", "line3"},
			maxLines: -1,
			want:     []string{"line1", "line2", "line3"},
		},
		{
			name:     "maxLines greater than length returns all lines",
			lines:    []string{"line1", "line2", "line3"},
			maxLines: 10,
			want:     []string{"line1", "line2", "line3"},
		},
		{
			name:     "maxLines equal to length returns all lines",
			lines:    []string{"line1", "line2", "line3"},
			maxLines: 3,
			want:     []string{"line1", "line2", "line3"},
		},
		{
			name:     "maxLines less than length returns limited lines",
			lines:    []string{"line1", "line2", "line3", "line4", "line5"},
			maxLines: 2,
			want:     []string{"line1", "line2"},
		},
		{
			name:     "maxLines 1 returns first line",
			lines:    []string{"line1", "line2", "line3"},
			maxLines: 1,
			want:     []string{"line1"},
		},
		{
			name:     "empty lines with maxLines",
			lines:    []string{},
			maxLines: 5,
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := limitLines(tt.lines, tt.maxLines)
			if len(got) != len(tt.want) {
				t.Errorf("limitLines() got %d lines, want %d lines", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("limitLines()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestValidateBridgeName(t *testing.T) {
	tests := []struct {
		name    string
		bridge  string
		wantErr bool
	}{
		{
			name:    "valid bridge name with hyphen",
			bridge:  "br-int",
			wantErr: false,
		},
		{
			name:    "valid bridge name with underscore",
			bridge:  "br_ex",
			wantErr: false,
		},
		{
			name:    "valid bridge name alphanumeric",
			bridge:  "br0",
			wantErr: false,
		},
		{
			name:    "valid bridge name mixed",
			bridge:  "br-local_123",
			wantErr: false,
		},
		{
			name:    "empty bridge name returns error",
			bridge:  "",
			wantErr: true,
		},
		{
			name:    "bridge name with space returns error",
			bridge:  "br int",
			wantErr: true,
		},
		{
			name:    "bridge name with slash returns error",
			bridge:  "br/int",
			wantErr: true,
		},
		{
			name:    "bridge name with semicolon returns error",
			bridge:  "br;int",
			wantErr: true,
		},
		{
			name:    "bridge name with pipe returns error",
			bridge:  "br|int",
			wantErr: true,
		},
		{
			name:    "bridge name with dollar returns error",
			bridge:  "br$int",
			wantErr: true,
		},
		{
			name:    "bridge name with backtick returns error",
			bridge:  "br`int",
			wantErr: true,
		},
		{
			name:    "bridge name with special chars returns error",
			bridge:  "br@int",
			wantErr: true,
		},
		{
			name:    "bridge name with parentheses returns error",
			bridge:  "br(int)",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBridgeName(tt.bridge)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBridgeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFlowSpec(t *testing.T) {
	tests := []struct {
		name    string
		flow    string
		wantErr bool
	}{
		{
			name:    "valid flow with in_port",
			flow:    "in_port=1,icmp",
			wantErr: false,
		},
		{
			name:    "valid flow with IP addresses",
			flow:    "in_port=2,ip,nw_src=192.168.1.10,nw_dst=192.168.1.20",
			wantErr: false,
		},
		{
			name:    "valid flow with TCP ports",
			flow:    "in_port=3,tcp,nw_src=10.0.0.1,nw_dst=10.0.0.2,tp_src=12345,tp_dst=80",
			wantErr: false,
		},
		{
			name:    "valid flow with brackets",
			flow:    "in_port=1,ip,nw_src=10.244.0.5,nw_dst=10.96.0.1",
			wantErr: false,
		},
		{
			name:    "valid flow with parentheses",
			flow:    "flow(test)",
			wantErr: false,
		},
		{
			name:    "valid flow with forward slash",
			flow:    "in_port=1,ip,nw_src=10.0.0.0/24",
			wantErr: false,
		},
		{
			name:    "empty flow returns error",
			flow:    "",
			wantErr: true,
		},
		{
			name:    "flow with semicolon returns error",
			flow:    "in_port=1;icmp",
			wantErr: true,
		},
		{
			name:    "flow with ampersand returns error",
			flow:    "in_port=1&icmp",
			wantErr: true,
		},
		{
			name:    "flow with pipe returns error",
			flow:    "in_port=1|icmp",
			wantErr: true,
		},
		{
			name:    "flow with dollar returns error",
			flow:    "in_port=$1,icmp",
			wantErr: true,
		},
		{
			name:    "flow with backtick returns error",
			flow:    "in_port=`1`,icmp",
			wantErr: true,
		},
		{
			name:    "flow with less than returns error",
			flow:    "in_port=1,icmp<test",
			wantErr: true,
		},
		{
			name:    "flow with greater than returns error",
			flow:    "in_port=1,icmp>test",
			wantErr: true,
		},
		{
			name:    "flow with backslash returns error",
			flow:    "in_port=1\\icmp",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlowSpec(tt.flow)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlowSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
