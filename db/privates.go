package db

import (
	"os"
	"path/filepath"
	"strings"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"errors"
)


// addError adds error to error list
func (d *Driver) addError(err error) *Driver {
	d.errors = append(d.errors, fmt.Errorf("simd: %v", err))
	return d
}

func (d *Driver) openDB(entity interface{}) ([]interface{}, error) {
	entityName, err:=d.getEntityName()
	
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

func(d *Driver) isDBOpened() bool {
	if(d.isOpened==false){
		err:=errors.New("should call Open() before doing any query on json file")
		d.addError(err)
	}
	return d.isOpened
}

func (d *Driver) getEntityName () (string, error) {
	typeName:=strings.Split(reflect.TypeOf(d.entityDealingWith).String(), ".")
	if len(typeName)<=0 {
		return "", fmt.Errorf("unable infer the type of the entity passed")
	}

	return typeName[len(typeName)-1], nil
}

func (d *Driver) readAppend(entity interface{}) (err error) {
	result, err:=d.openDB(entity)
	mergedArray, err:=mergeToExisting(result, entity)
	err=d.writeAll(mergedArray)
	return
}

func(d *Driver) writeAll(entities []interface{}) (err error) {
	entityName, err:=d.getEntityName()
	file:=filepath.Join(d.dir, entityName)
	f, err:=os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	b, err:=json.MarshalIndent(entities, "", "\t")
	if err!=nil {
		return
	}
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
				return nil, fmt.Errorf("invalid operator %s " + q.operator)
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
