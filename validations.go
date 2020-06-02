package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/2HgO/gin-web-server/configs/db"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(validationFieldName)
		v.RegisterValidation("movie", dbValidation(db.Movie))
		v.RegisterValidation("category", dbValidation(db.Category))
		log.Println("Succesfully registered validations")
	}
}

func validationFieldName(fld reflect.StructField) string {
	var name string
	if tag, ok := fld.Tag.Lookup("form"); ok {
		name = strings.SplitN(tag, ",", 2)[0]
	} else {
		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	}
	return name
}

func dbValidation(c *db.Collection) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Interface()
		exists, err := c.Exists(map[string]interface{}{"_id": val})
		return exists && err == nil
	}
}
