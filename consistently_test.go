package consistently_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ftomza/consistently"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name string
	consistently.Consistently
}

type UserCustom struct {
	gorm.Model
	Name     string
	Version2 string `consistently:"version"`
}

func TestGoConsistently(t *testing.T) {

	user := User{Name: "FtomZa"}

	db.Save(&user)

	if user.Version == "" {
		t.Errorf("Must contain a version of the data")
	}

	user.Name = "gopa"

	beforeUpdateVesrion := user.Version

	errs := db.Save(&user).GetErrors()

	errFound := false
	for _, v := range errs {

		if v == consistently.ErrVersionNotValid {
			errFound = true
		}
	}

	if errFound {
		t.Errorf("Must not contain ErrVersionNotValid")
	}

	if user.Name != "gopa" {
		t.Errorf("Must equal user.Name")
	}

	userCheck := User{}

	db.First(&userCheck, user.ID)

	if beforeUpdateVesrion == userCheck.Version {
		t.Errorf("Must not equal new version and old version")
	}

	user.Version = "bad"

	errs = db.Save(&user).GetErrors()

	errFound = false
	for _, v := range errs {

		if v == consistently.ErrVersionNotValid {
			errFound = true
		}
	}

	if !errFound {
		t.Errorf("Must contain ErrVersionNotValid")
	}
}

func TestGoConsistentlyCustomField(t *testing.T) {

	user := UserCustom{Name: "FtomZa"}

	db.Save(&user)

	if user.Version2 == "" {
		t.Errorf("Must contain a version of the data")
	}

	beforeUpdateVesrion := user.Version2

	errs := db.Model(&user).Update(&user).GetErrors()

	errFound := false
	for _, v := range errs {

		if v == consistently.ErrVersionNotValid {
			errFound = true
		}
	}

	if errFound {
		t.Errorf("Must not contain ErrVersionNotValid")
	}

	userCheck := UserCustom{}

	db.First(&userCheck, user.ID)

	if beforeUpdateVesrion == userCheck.Version2 {
		t.Errorf("Must not equal new version and old version")
	}

	user.Version2 = "bad"

	errs = db.Model(&user).Update(&user).GetErrors()

	errFound = false
	for _, v := range errs {

		if v == consistently.ErrVersionNotValid {
			errFound = true
		}
	}

	if !errFound {
		t.Errorf("Must contain ErrVersionNotValid")
	}
}

func TestGoConsistentlyAsync(t *testing.T) {

	user := User{Name: "FtomZa"}

	db.Save(&user)

	var wg sync.WaitGroup
	wg.Add(2)

	handler := func(fail bool, dur time.Duration) {
		defer wg.Done()

		testUser := User{}
		db.Find(&testUser, user.ID)

		time.Sleep(dur * time.Second)

		errs := db.Model(&testUser).Update(testUser).GetErrors()

		errFound := false
		for _, v := range errs {

			if v == consistently.ErrVersionNotValid {
				errFound = true
			}
		}

		if errFound && !fail {
			t.Errorf("Must not contain ErrVersionNotValid")
		}

		if !errFound && fail {
			t.Errorf("Must contain ErrVersionNotValid")
		}
	}

	go handler(false, 2)
	go handler(true, 4)

	wg.Wait()
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	consistently.RegisterCallbacks(db)
	tables := []interface{}{&User{}, &UserCustom{}}
	for _, table := range tables {
		if err := db.DropTableIfExists(table).Error; err != nil {
			panic(err)
		}
		db.AutoMigrate(table)
	}
}
