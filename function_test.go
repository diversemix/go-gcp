// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"testing"

	visionapi "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// not nice but works
var _ = (func() interface{} {
	_testing = true
	return nil
}())

func TestDetectAndCrop(t *testing.T) {
	type args struct {
		ctx context.Context
		e   GCSEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DetectAndCrop(tt.args.ctx, tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("DetectAndCrop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_crop(t *testing.T) {
	type args struct {
		ctx          context.Context
		inputBucket  string
		outputBucket string
		inputName    string
		outputName   string
		minX         int32
		minY         int32
		maxX         int32
		maxY         int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := crop(tt.args.ctx, tt.args.inputBucket, tt.args.outputBucket, tt.args.inputName, tt.args.outputName, tt.args.minX, tt.args.minY, tt.args.maxX, tt.args.maxY); (err != nil) != tt.wantErr {
				t.Errorf("crop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findBoundingRect(t *testing.T) {
	type args struct {
		vertices []*visionapi.Vertex
	}
	v1 := visionapi.Vertex{X: 0, Y: 0}
	v2 := visionapi.Vertex{X: 0, Y: 0}
	var vertices []*visionapi.Vertex
	vertices = append(vertices, &v1)
	vertices = append(vertices, &v2)
	var testStruct args
	testStruct.vertices = vertices

	tests := []struct {
		name  string
		args  args
		want  int32
		want1 int32
		want2 int32
		want3 int32
	}{
		{name: "fred", args: testStruct, want: 0, want1: 0, want2: 0, want3: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := findBoundingRect(tt.args.vertices)
			if got != tt.want {
				t.Errorf("findBoundingRect() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findBoundingRect() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("findBoundingRect() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("findBoundingRect() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
