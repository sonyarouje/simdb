package db

import (
	"errors"
)

//Entity any structure wanted to persist to json should implement this interface
//ID return the id value and field name that stores the id
/*e.g 
	type Customer struct {
		CustID string `json:"cust_id"`
		Address string `json:"address"`
	}

	func (c Customer) ID() (value interface{}, jsonField string) {
		value=c.CustID
		jsonField="cust_id"
		return
	}
*/
type Entity interface {
	ID() (value interface{}, jsonField string)
}

// empty represents an empty result
var empty interface{}

// query describes a query
type query struct {
	key, operator string
	value         interface{}
}

type Driver struct {
	dir string
	queries         [][]query            // nested queries
	queryIndex      int
	queryMap        map[string]QueryFunc // contains query functions
	jsonContent     interface{}          // copy of original decoded json data for further processing
	errors          []error              // contains all the errors when processing
	originalJSON	interface{}			 //actual json when opening the json file
	isOpened		bool
}

//New creates a new database driver. Pass the directory to store the db files.
func New(dir string) (*Driver, error) {
	driver:= &Driver {
		dir:dir,
		queryMap: loadDefaultQueryMap(),
	}
	err:= createDirIfNotExist(dir)
	return driver, err
}

//Open will open the json file based on the entity passed.
//Once the file is open you can apply where conditions or get operation.
func (d *Driver) Open(entity interface{}) *Driver {
	db, err:=d.openDB(entity)
	d.originalJSON=db
	d.isOpened=true
	if(err!=nil){
		d.addError(err)
	}
	return d
}

//Errors will return errors encountered while performing any operations
func (d * Driver) Errors () []error {
	return d.errors
}

//Insert the data to the json db. Insert will identify the type of the 
//entity and insert the entity to the specific json file.
//If the file not exist then will create a new file
func (d *Driver) Insert(entity Entity) (err error) {
	err=d.readAppend(entity)
	return
}

// Where builds a where clause. e.g: Where("name", "=", "doe")
func (d *Driver) Where(key, cond string, val interface{}) *Driver {
	q := query{
		key:      key,
		operator: cond,
		value:    val,
	}
	if d.queryIndex == 0 && len(d.queries) == 0 {
		qq := []query{}
		qq = append(qq, q)
		d.queries = append(d.queries, qq)
	} else {
		d.queries[d.queryIndex] = append(d.queries[d.queryIndex], q)
	}

	return d
}

//Get the result from the json db. If no where condition then return all the data from json
func(d *Driver) Get() []interface{}{
	if(d.isOpened==false){
		err:=errors.New("should call Open() before doing any query on json file")
		d.addError(err)
		return nil
	}
	if len(d.queries) > 0 {
		d.processQuery()
	}else{
		d.jsonContent=d.originalJSON
	}
	d.queryIndex = 0
	if aa, ok := d.jsonContent.([]interface{}); ok {
		return aa
	}
	return nil
}
