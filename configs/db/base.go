package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocBase struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt *time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type SoftDocBase struct {
	DocBase `json:",inline" bson:",inline"`
}

func (d *DocBase) GetID() *primitive.ObjectID {
	return d.ID
}

func (d *DocBase) SetID(id interface{}) {
	if d == nil {
		panic("Cannot set ID of nil pointer")
	}
	if i, ok := id.(primitive.ObjectID); ok {
		d.ID = &i
		return
	}
	panic("Invalid ID")
}

func (d *DocBase) SetCreated(t time.Time) {
	if d == nil {
		panic("Cannot set created of nil pointer")
	}
	d.CreatedAt = &t
}

func (d *DocBase) SetUpdated(t time.Time) {
	if d == nil {
		panic("Cannot set updated of nil pointer")
	}
	d.UpdatedAt = &t
}

func (s *SoftDocBase) SoftDelete() {}
