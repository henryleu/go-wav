package wav

import (
	"testing"
)

func TestNewReaderFromFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new reader from an existent file",
			args: args{
				fileName: "./testdata/a.wav",
			},
			wantErr: false,
		},
		{
			name: "new reader from an non-existent file",
			args: args{
				fileName: "./testdata/nofile.wav",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewReaderFromFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReaderFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
