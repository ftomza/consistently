# Consistently

[![GoDoc](https://godoc.org/github.com/ftomza/consistently?status.svg)](https://godoc.org/github.com/ftomza/consistently)
[![Go Report Card](https://goreportcard.com/badge/github.com/ftomza/consistently)](https://goreportcard.com/report/github.com/ftomza/consistently)
[![Build Status](https://travis-ci.org/ftomza/consistently.svg?branch=master)](https://travis-ci.org/ftomza/consistently)
[![Coverage Status](https://coveralls.io/repos/github/ftomza/consistently/badge.svg)](https://coveralls.io/github/ftomza/consistently)

Plugin for [GORM](https://github.com/jinzhu/gorm) checking read data before updating for consistency

### Register GORM Callbacks
Consistently uses [GORM](https://github.com/jinzhu/gorm) callbacks to handle *gorm:update* and *gorm:create*, so you will need to register callbacks first:

```go
import (
  "github.com/jinzhu/gorm"
  "github.com/ftomza/consistently
)

func main() {
  db, err := gorm.Open("sqlite3", "demo_db")

  consistently.RegisterCallbacks(db)
}
```

### Usage

After callbacks are registered, an attempt to create or update any record checks all fields of the structure for the tag *`consistently:"version"`*. If such a field is found, a new data version value will be assigned when the record is created, if it is updated, the version of the recorded data with the version in the database will be checked.

```go
type User struct {
	gorm.Model
	Name string
	consistently.Consistently
}


func main() {
	
  db, err := gorm.Open("sqlite3", "demo_db")

  consistently.RegisterCallbacks(db)
  
  user := User{Name: "FtomZa"}
  
  db.Save(&user)
  
  fmt.printf("Version: %s\n", user.Version)
  
  user.Name = "FtomZa2"
  
  db.Update(&user)
  
  fmt.printf("Version: %s\n", user.Version)
}
```

## Custom field for save version

If it is not possible to store version data in the *version* field by default, use the *`consistently:"version"`* tag to mark the field as a data version.

```go
type User struct {
	gorm.Model
	Name string
	MyVersion string `consistently:"version"`
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
