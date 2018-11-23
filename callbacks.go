package consistently

import (
	"github.com/jinzhu/gorm"
)

const (
	tagKeyConsistently = "consistently"
	tagValueVersion    = "version"
)

func callbackConsistently(scope *gorm.Scope) {

	if _, ok := scope.Get("gorm:update_column"); ok {
		return
	}

	if scope.HasError() {
		return
	}

	if scope.Value == nil {
		return
	}

	for _, field := range scope.Fields() {

		if val, ok := field.Tag.Lookup(tagKeyConsistently); !(ok && val == tagValueVersion) {
			continue
		}

		saveVersion := field.Field.String()

		if saveVersion != "" {
			currentValue := scope.New(scope.Value)

			scope.DB().First(currentValue.Value)

			if scope.HasError() {
				return
			}

			currentField, _ := currentValue.FieldByName(field.Name)
			currentVersion := currentField.Field.String()
			if currentVersion != saveVersion {
				scope.Err(ErrVersionNotValid)
				return
			}
		}

		scope.SetColumn(field.Name, string(randASCIIBytes(20)))
		return
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {

	callback := db.Callback()

	if callback.Create().Get("consistently:before_create") == nil {
		callback.Create().Before("gorm:before_create").Register("consistently:before_create", callbackConsistently)
	}

	if callback.Update().Get("consistently:before_update") == nil {
		callback.Update().Before("gorm:before_update").Register("consistently:before_update", callbackConsistently)
	}
}
