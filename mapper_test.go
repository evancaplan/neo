package sbdb_cad

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestSBDecoder_Decode(t *testing.T) {
	type args struct {
		input  interface{}
		output interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error occurs in mapstructure.Decode",
			args: args{
				input:  make(map[string]string),
				output: SbCAD{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &SbCADDecoder{}
			if err := sd.Decode(tt.args.input, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSBMapper_Map(t *testing.T) {
	tests := []struct {
		name    string
		decoder Decoder
		args    *http.Response
		sbr     SbCADResponse
		want    []SbCAD
		wantErr bool
	}{
		{
			name:    "Returns Array of SBs with no error",
			decoder: &MockDecoder{err: nil},
			args:    &http.Response{},
			sbr: SbCADResponse{
				Count:  "1",
				Fields: []string{},
				Data:   [][]string{{}},
			},
			want:    []SbCAD{{}},
			wantErr: false,
		}, {
			name:    "Zero count error occurs",
			decoder: &MockDecoder{err: nil},
			args:    &http.Response{},
			sbr: SbCADResponse{
				Count:  "0",
				Fields: []string{},
				Data:   [][]string{{}},
			},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Error Occurs in decoder",
			decoder: &MockDecoder{err: errors.New("error in decoder")},
			args:    &http.Response{},
			sbr: SbCADResponse{
				Count:  "1",
				Fields: []string{},
				Data:   [][]string{{}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SbCADMapper{
				Decoder: tt.decoder,
			}

			respBytes, err := json.Marshal(tt.sbr)
			if err != nil {
				t.Errorf("Error marshalling response bytes")
			}
			tt.args.Body = ioutil.NopCloser(bytes.NewBuffer(respBytes))

			got, err := sb.Map(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSBMapper_mapSBResArrayToStruct(t *testing.T) {
	type args struct {
		res    []string
		fields []string
	}
	tests := []struct {
		name    string
		decoder Decoder
		args    args
		want    *SbCAD
		wantErr bool
	}{
		{
			name:    "SbCAD Response Array mapped to Struct, no error",
			decoder: &MockDecoder{err: nil},
			args: args{
				res:    []string{},
				fields: []string{},
			},
			want: &SbCAD{},
		}, {
			name:    "Error occurs in decoder",
			decoder: &MockDecoder{err: errors.New("error occurred in decoder")},
			args: args{
				res:    []string{},
				fields: []string{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SbCADMapper{
				Decoder: tt.decoder,
			}
			got, err := sb.mapSbCADResArrayToStruct(tt.args.res, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapSbCADResArrayToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapSbCADResArrayToStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSBMapper_mapSBResToSB(t *testing.T) {
	tests := []struct {
		name       string
		decoder    Decoder
		sbResponse *SbCADResponse
		want       []SbCAD
		wantErr    bool
	}{
		{
			name:    "Return SbCAD, no error occurs",
			decoder: &MockDecoder{err: nil},
			sbResponse: &SbCADResponse{
				Count:  "1",
				Fields: []string{},
				Data:   [][]string{{}},
			},
			want:    []SbCAD{{}},
			wantErr: false,
		},
		{
			name:
			"Error occurs in decoder",
			decoder: &MockDecoder{err: errors.New("error decoding")},
			sbResponse: &SbCADResponse{
				Count:  "0",
				Fields: []string{},
				Data:   [][]string{{}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SbCADMapper{
				Decoder: tt.decoder,
			}
			got, err := sb.mapSBCADResToSBCAD(tt.sbResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapNeoResToNeo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapNeoResToNeo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockDecoder struct {
	err error
}

func (md *MockDecoder) Decode(input interface{}, output interface{}) error {
	return md.err
}
