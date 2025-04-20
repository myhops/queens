package pkg

import (
	"reflect"
	"testing"
)

func TestFields_Row(t *testing.T) {
	type fields struct {
		size int
	}
	type args struct {
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Position
	}{
		{
			name:   "Row 5",
			fields: fields{size: 5},
			args: args{
				y: 5,
			},
			want: []Position{
				{x: 0, y: 5},
				{x: 1, y: 5},
				{x: 2, y: 5},
				{x: 3, y: 5},
				{x: 4, y: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fields{
				size: tt.fields.size,
			}
			if got := f.Row(tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields.Row() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFields_Column(t *testing.T) {
	type fields struct {
		size int
	}
	type args struct {
		x int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Position
	}{
		{
			name:   "Column 5",
			fields: fields{size: 5},
			args: args{
				x: 5,
			},
			want: []Position{
				{x: 5, y: 0},
				{x: 5, y: 1},
				{x: 5, y: 2},
				{x: 5, y: 3},
				{x: 5, y: 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fields{
				size: tt.fields.size,
			}
			if got := f.Column(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields.Column() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFields_Around(t *testing.T) {
	type fields struct {
		size int
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Position
	}{
		// TODO: Add test cases.
		{
			name: "center",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 2,
				y: 2,
			},
			want: []Position{
				{1, 1}, {1, 2}, {1, 3},
				{2, 1}, {2, 3},
				{3, 1}, {3, 2}, {3, 3},
			},
		},
		{
			name: "middle top",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 2,
				y: 0,
			},
			want: []Position{
				{1, 0}, {1, 1},
				{2, 1},
				{3, 0}, {3, 1},
			},
		},
		{
			name: "middle bottom",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 2,
				y: 4,
			},
			want: []Position{
				{1, 3}, {1, 4},
				{2, 3},
				{3, 3}, {3, 4},
			},
		},
		{
			name: "middle left",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 0,
				y: 2,
			},
			want: []Position{
				{0, 1}, {0, 3},
				{1, 1}, {1, 2}, {1, 3},
			},
		},
		{
			name: "middle right",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 4,
				y: 2,
			},
			want: []Position{
				{3, 1}, {3, 2}, {3, 3},
				{4, 1}, {4, 3},
			},
		},
		{
			name: "top left",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 0,
				y: 0,
			},
			want: []Position{
				{0, 1},
				{1, 0}, {1, 1},
			},
		},
		{
			name: "bottom right",
			fields: fields{
				size: 5,
			},
			args: args{
				x: 4,
				y: 4,
			},
			want: []Position{
				{3, 3}, {3, 4},
				{4, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fields{
				size: tt.fields.size,
			}
			if got := f.Around(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields.Around() = %v, want %v", got, tt.want)
			}
		})
	}
}
