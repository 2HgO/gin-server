package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	
	"github.com/2HgO/gin-web-server/errors"
)

type Date struct {
	time.Time
}

func NewDate(t time.Time) *Date {
	return &Date{t}
}

func (d *Date) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	stamp, err := strconv.ParseInt(strInput, 10, 64)
	if err != nil {
		t, err := time.Parse("2006-01-02", strInput)
		if err != nil {
			return errors.ValidationError
		}
		d.Time = t.In(time.UTC)
		return nil
	}

	newTime := time.Unix(stamp/1000, 0)
	d.Time = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), 0, 0, 0, 0, time.UTC)
	return nil
}

func (d *Date) UnmarshalBSON(data []byte) error {
	val := bson.RawValue{
		Type:  bsontype.DateTime,
		Value: data,
	}
	stamp, ok := val.DateTimeOK()
	if !ok {
		return errors.ValidationError
	}

	newTime := time.Unix(stamp/1000, 0)
	d.Time = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), 0, 0, 0, 0, time.UTC)
	return nil
}

func (d *Date) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(d.Time)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}
