package utils

import "testing"

func TestDataHandler_Start(t *testing.T) {
	type fields struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				filePath: "../../data/sougouCA_seg_simple.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewDataHandler(tt.fields.filePath)
			if err := h.Start(nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
