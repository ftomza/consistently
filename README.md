# Consistently
Plugin for [GORM](https://github.com/jinzhu/gorm) checking read data before updating for consistency

### Register GORM Callbacks
Validations uses [GORM](https://github.com/jinzhu/gorm) callbacks to handle *gorm:update* and *gorm:create*, so you will need to register callbacks first:

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

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
