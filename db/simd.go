package db

import (
	"os"
	"path/filepath"
	"strings"
	"errors"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"fmt"
)
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

//New creates a new database driver
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

//Errors will return errors encounters while performing any operation
func (d * Driver) Errors () []error {
	return d.errors
}

//Insert the data to the json db. Insert will identify the type of the 
//entity and create/append the entity to the specific json file.
func (d *Driver) Insert(entity interface{}) error {
	err:=d.readAppend(entity)
	return err
}


// Where builds a where clause. e.g: Where("name", "contains", "doe")
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

// addError adds error to error list
func (d *Driver) addError(err error) *Driver {
	d.errors = append(d.errors, fmt.Errorf("simd: %v", err))
	return d
}

func (d *Driver) openDB(entity interface{}) ([]interface{}, error) {
	entityName, err:=d.getEntityName(entity)
	if(err!=nil){
		return nil, err
	}
	file:=filepath.Join(d.dir, entityName)

	f, err:=os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	if(err!=nil){
		return nil, err
	}

	b, readErr:=ioutil.ReadFile(file)
	if readErr!=nil {
		return nil, readErr
	}
	array:=make([]interface{}, 0)
	json.Unmarshal(b,&array)

	return array, nil
}

func (d *Driver) getEntityName (entity interface{}) (string, error) {
	typeName:=strings.Split(reflect.TypeOf(entity).String(), ".")
	if len(typeName)<=0 {
		err:= errors.New("unable infer the type of the entity passed")
		return "", err
	}

	return typeName[len(typeName)-1], nil
}

func (d *Driver) readAppend(entity interface{}) (err error) {
	entityName, err:=d.getEntityName(entity)
	file:=filepath.Join(d.dir, entityName)
	f, err:=os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	result, err:=d.openDB(entity)
	b, err:=mergeToExisting(result, entity)

	f.Truncate(0)
	f.Seek(0,0)
	f.Write(b)
	f.Sync()

	return
}

// findInArray traverses through a list and returns the value list.
// This helps to process Where/OrWhere queries
func (d *Driver) findInArray(aa []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, a := range aa {
		if m, ok := a.(map[string]interface{}); ok {
			findResult, err:=d.findInMap(m)
			if(err==nil){
				result = append(result, findResult...)
			}else{
				d.addError(err)
			}
		}
	}
	return result
}

// findInMap traverses through a map and returns the matched value list.
// This helps to process Where/OrWhere queries
func (d *Driver) findInMap(vm map[string]interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0)
	orPassed := false
	for _, qList := range d.queries {
		andPassed := true
		for _, q := range qList {
			cf, ok := d.queryMap[q.operator]
			if !ok {
				err:=errors.New ("invalid operator %s " + q.operator)
				return nil, err
			}
			nv, errnv := getNestedValue(vm, q.key)
			if errnv != nil {
				return nil, errnv
			} else {
				qb, err := cf(nv, q.value)
				if err != nil {
					return nil, err
				}
				andPassed = andPassed && qb
			}
		}
		orPassed = orPassed || andPassed
	}
	if orPassed {
		result = append(result, vm)
	}
	return result, nil
}

// processQuery makes the result
func (d *Driver) processQuery() *Driver {
	if aa, ok := d.originalJSON.([]interface{}); ok {
		d.jsonContent = d.findInArray(aa)
	}
	return d
}
