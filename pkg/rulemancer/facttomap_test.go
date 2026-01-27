package rulemancer

import (
	"strings"
	"testing"
)

func TestGenericFactToMap(t *testing.T) {
	tests := []struct {
		name        string
		statusItem  string
		factList    string
		want        []map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:       "single statusItem with multiple key-value pairs",
			statusItem: "health",
			factList:   "(health (current 100) (max 150))",
			want: []map[string]string{
				{
					"current": "100",
					"max":     "150",
				},
			},
			wantErr: false,
		},
		{
			name:       "multiple statusItems with same fields",
			statusItem: "health",
			factList:   "(health (current 100) (max 150)) (health (current 80) (max 120))",
			want: []map[string]string{
				{
					"current": "100",
					"max":     "150",
				},
				{
					"current": "80",
					"max":     "120",
				},
			},
			wantErr: false,
		},
		{
			name:       "single key-value pair",
			statusItem: "player",
			factList:   "(player (name John))",
			want: []map[string]string{
				{
					"name": "John",
				},
			},
			wantErr: false,
		},
		{
			name:       "statusItem with extra whitespace",
			statusItem: "stats",
			factList:   "(stats  (level 5)  (xp 1000) )",
			want: []map[string]string{
				{
					"level": "5",
					"xp":    "1000",
				},
			},
			wantErr: false,
		},
		{
			name:       "no matching statusItem",
			statusItem: "health",
			factList:   "(mana (current 50) (max 100))",
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "empty factList",
			statusItem: "health",
			factList:   "",
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "statusItem exists but no key-value pairs",
			statusItem: "empty",
			factList:   "(empty )",
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "multiple different statusItems in list",
			statusItem: "position",
			factList:   "(health (current 100)) (position (x 10) (y 20)) (mana (current 50))",
			want: []map[string]string{
				{
					"x": "10",
					"y": "20",
				},
			},
			wantErr: false,
		},
		{
			name:       "values with special characters",
			statusItem: "item",
			factList:   "(item (type sword) (name Fire-Blade))",
			want: []map[string]string{
				{
					"type": "sword",
					"name": "Fire-Blade",
				},
			},
			wantErr: false,
		},
		{
			name:       "numeric values",
			statusItem: "coords",
			factList:   "(coords (x 123.45) (y -67.89))",
			want: []map[string]string{
				{
					"x": "123.45",
					"y": "-67.89",
				},
			},
			wantErr: false,
		},
		{
			name:       "multiple items with consistent fields",
			statusItem: "enemy",
			factList:   "(enemy (type orc) (hp 50)) (enemy (type goblin) (hp 30)) (enemy (type troll) (hp 100))",
			want: []map[string]string{
				{
					"type": "orc",
					"hp":   "50",
				},
				{
					"type": "goblin",
					"hp":   "30",
				},
				{
					"type": "troll",
					"hp":   "100",
				},
			},
			wantErr: false,
		},
		{
			name:       "statusItem name with special chars that need escaping",
			statusItem: "stat.health",
			factList:   "(stat.health (current 100) (max 200))",
			want: []map[string]string{
				{
					"current": "100",
					"max":     "200",
				},
			},
			wantErr: false,
		},
		{
			name:        "inconsistent fields - different number of fields",
			statusItem:  "player",
			factList:    "(player (name John) (level 5)) (player (name Jane))",
			want:        nil,
			wantErr:     true,
			errContains: "inconsistent fields",
		},
		{
			name:        "inconsistent fields - different field names",
			statusItem:  "player",
			factList:    "(player (name John) (level 5)) (player (name Jane) (rank 3))",
			want:        nil,
			wantErr:     true,
			errContains: "inconsistent fields",
		},
		{
			name:        "inconsistent fields - missing field",
			statusItem:  "item",
			factList:    "(item (id 1) (name sword) (type weapon)) (item (id 2) (type weapon))",
			want:        nil,
			wantErr:     true,
			errContains: "inconsistent fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genericFactToMap(nil, tt.statusItem, tt.factList)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenericFactToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GenericFactToMap() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GenericFactToMap() returned %d items, want %d items", len(got), len(tt.want))
					t.Errorf("Got: %v", got)
					t.Errorf("Want: %v", tt.want)
					return
				}

				for i, wantMap := range tt.want {
					gotMap := got[i]
					if len(gotMap) != len(wantMap) {
						t.Errorf("GenericFactToMap() item %d has %d fields, want %d fields", i, len(gotMap), len(wantMap))
						return
					}
					for key, wantValue := range wantMap {
						gotValue, exists := gotMap[key]
						if !exists {
							t.Errorf("GenericFactToMap() item %d missing key %q", i, key)
						} else if gotValue != wantValue {
							t.Errorf("GenericFactToMap() item %d for key %q = %v, want %v", i, key, gotValue, wantValue)
						}
					}
				}
			}
		})
	}
}
