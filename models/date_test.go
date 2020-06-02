package models

import (
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/2HgO/gin-web-server/errors"
)

func TestNewDate(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "Ensure correct Date object is returned",
			args: args{
				date: time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC),
			},
			want: Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := NewDate(test.args.date)
			Equal(t, test.want.Time, d.Time)
		})
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	type args struct {
		datestring string
	}
	tests := []struct {
		name string
		args args
		want Date
		err  error
	}{
		{
			name: "Ensure correct unmarshalling of ms timestamp",
			args: args{
				datestring: "1581984000000",
			},
			want: Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "Ensure correct unmarshalling of Date string",
			args: args{
				datestring: "2020-02-18",
			},
			want: Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "Ensure error is returned for incorrect Date format",
			args: args{
				datestring: "02-18-2020",
			},
			err: errors.ValidationError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := new(Date)
			err := d.UnmarshalJSON([]byte(test.args.datestring))
			Equal(t, test.want.Time, d.Time)
			Equal(t, test.err, err)
		})
	}
}

func TestDateMarshalJSON(t *testing.T) {
	type args struct {
		d Date
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "Ensure correct datetime string in marshalled json-encoded string",
			args: args{
				d: Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
			},
			want: `"2020-02-18"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := test.args.d.MarshalJSON()
			Equal(t, test.want, string(bytes))
			Equal(t, test.err, err)
		})
	}
}

func TestDateMarshalBSONValue(t *testing.T) {
	type args struct {
		time time.Time
		d    Date
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Ensure correct datetime string in marshalled bson-encoded string",
			args: args{
				time: time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC),
				d:    Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedType, expectedVal, expectedErr := bson.MarshalValue(test.args.time)
			bsonType, bytes, err := test.args.d.MarshalBSONValue()
			Equal(t, expectedVal, bytes)
			Equal(t, expectedErr, err)
			Equal(t, expectedType, bsonType)
		})
	}
}

func TestDateUnmarshalBSON(t *testing.T) {
	type args struct {
		time []byte
	}
	_, bytes, _ := bson.MarshalValue(time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC))
	tests := []struct {
		name string
		args args
		want Date
		err  error
	}{
		{
			name: "Ensure correct unmarshalling of bson datetime",
			args: args{
				time: bytes,
			},
			want: Date{time.Date(2020, 2, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "Ensure error is returned for incorrect bson datetime format",
			args: args{
				time: []byte(`"z"`),
			},
			err: errors.ValidationError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := new(Date)
			err := d.UnmarshalBSON(test.args.time)
			Equal(t, test.want.Time, d.Time)
			Equal(t, test.err, err)
		})
	}
}
