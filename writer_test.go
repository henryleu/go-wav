package wav

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func printFmtChunk(chk *FmtChunk) {
	fmt.Printf("FmtChunk.Data.Channel:\t%v\n", chk.Data.Channel)
	fmt.Printf("FmtChunk.Data.SamplesPerSec:\t%v\n", chk.Data.SamplesPerSec)
	fmt.Printf("FmtChunk.Data.BytesPerSec:\t%v\n", chk.Data.BytesPerSec)
	fmt.Printf("FmtChunk.Data.BlockSize:\t%v\n", chk.Data.BlockSize)
	fmt.Printf("FmtChunk.Data.BitsPerSamples:\t%v\n", chk.Data.BitsPerSamples)
}

func printWriterParam(p *WriterParam) {
	fmt.Printf("Channel:\t%v\n", p.Channel)
	fmt.Printf("SampleRate:\t%v\n", p.SampleRate)
	fmt.Printf("BitsPerSample:\t%v\n", p.BitsPerSample)
}

func TestNewWriter_and_write(t *testing.T) {
	type args struct {
		param WriterParam
		from  string
		to    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Writer
		wantErr bool
	}{
		{
			name: "read and write to a temp file",
			args: args{
				param: WriterParam{
					Channel:       1,
					SampleRate:    8000,
					BitsPerSample: 16,
				},
				from: "./testdata/a.wav",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "*.wav")
			if err != nil {
				t.Error(err)
			}
			tt.args.param.Out = file
			tt.args.to = file.Name()
			defer os.Remove(file.Name())
			fmt.Println(tt.args.to)
			printWriterParam(&tt.args.param)
			w, err := NewWriter(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r, err := NewReaderFromFile(tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReaderFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			printFmtChunk(r.FmtChunk)
			frame := make([]byte, r.FmtChunk.Data.BlockSize)
			for {
				_, err := io.ReadFull(r, frame)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					fmt.Println(tt.args.from, "EOF")
					break
				} else if err != nil {
					t.Errorf("io.ReadFull() error = %v", err)
				}
				_, err = w.Write(frame)
				if err != nil {
					t.Errorf("Fail to write wave data, error: %v", err)
				}
			}
			w.Close()

			nr, err := NewReaderFromFile(tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReaderFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			printFmtChunk(nr.FmtChunk)
			if !reflect.DeepEqual(r.FmtChunk, nr.FmtChunk) {
				t.Errorf("NewReaderFromFile() = %v, want %v", nr.FmtChunk, r.FmtChunk)
			}

		})
	}
}
