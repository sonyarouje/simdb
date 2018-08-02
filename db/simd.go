package db

import (
	"fmt"
	"errors"
	"github.com/mitchellh/mapstructure"
	// "reflect"
)

//Entity any structure wanted to persist to json should implement this interface.
//ID and Field will be used while doing update or delete operation.
//ID return the id value and field name that stores the id
/*e.g 
	type Customer struct {
		CustID string `json:"custid"`
		Name string `json:"name"`
		Address string `json:"address"`
	}

	func (c Customer) ID() (jsonField string, value interface{}) {
		value=c.CustID
		jsonField="custid"
		return
	}
*/
type Entity interface {
	ID() (jsonField string, value interface{})
}

// empty represents an empty result
var empty interface{}

// query describes a query
type query struct {
	key, operator string
	value         interface{}
}

//Driver contains all the state of db.
type Driver struct {
	dir string							 //directory name to store the db
	queries         [][]query            // nested queries
	queryIndex      int
	queryMap        map[string]QueryFunc // contains query functions
	jsonContent     interface{}          // copy of original decoded json data for further processing
	errors          []error              // contains all the errors when processing
	originalJSON	interface{}			 // actual json when opening the json file
	isOpened		bool
	entityDealingWith interface{}		 // keeps the entity the driver is dealing with, field will maintain only the last entity inserted or updated or opened
}

//New creates a new database driver. Accepts the directory name to store the db files.
//If the passed directory not exist then will create it.
func New(dir string) (*Driver, error) {
	driver:= &Driver {
		dir:dir,
		queryMap: loadDefaultQueryMap(),
	}
	err:= createDirIfNotExist(dir)
	return driver, err
}

//Open will open the json file db based on the entity passed.
//Once the file is open you can apply where conditions or get operation.
func (d *Driver) Open(entity interface{}) *Driver {
	d.entityDealingWith=entity

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

//Insert the entity to the json db. Insert will identify the type of the 
//entity and insert the entity to the specific json file basee on the type of the entity.
//If the file not exist then will create a new file
func (d *Driver) Insert(entity Entity) (err error) {
	d.entityDealingWith=entity
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

//Get the result from the json db as an array. If no where condition then return all the data from json
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

//First return the first record matching the condtion.
func(d *Driver) First() interface{} {
	records:=d.Get()
	
	if len(records)>0 {
		return records[0]
	}

	return nil
}

//ToEntity will converts the map to the passed structure.
//result parameter takes the result returned by Get() or First()
//out will take pointer to structure.
//e.g. 
// struct custOut Customer
// driver.ToEntity(result, &custOut)
// this function will fill the custOut with the values from the map
func (d *Driver) ToEntity(result interface{}, out interface{}) interface {}{

	err:=mapstructure.Decode(result, out)
	if(err!=nil) {
		panic(err)
	}
	// fmt.Printf("%#v \n", *tmp)
	return out
}

// func (d *Driver) ToEntityArray(result []interface{}, out interface{}) []interface{} {
// 	outArray:=make([]interface{}, 0)
// 	for _, item:=range result {
// 		// structType:=reflect.TypeOf(out)
// 		// structValue:=reflect.Zero(structType)
// 		// structInterface:=structValue.Interface()
// 		// newStruct:=structInterface
// 		fmt.Printf("%#v", item)
// 		fmt.Println("")
// 		tmp:=&out
// 		err:=mapstructure.Decode(item, tmp)
// 		if(err!=nil){
// 			panic(err)
// 		}

// 		fmt.Printf("%#v", out)
// 		fmt.Println("")
// 		outArray=append(outArray,out)
// 	}
// 	return outArray
// }

//Update the json data based on the id field/value pair
func (d *Driver) Update(entity Entity) (err error) {
	d.entityDealingWith=entity
	field, entityID:=entity.ID()
	records:= d.Open(entity).Get()
	couldUpdate:=false
	entName,_:=d.getEntityName()

	if(len(records)>0){
		for indx,item:= range records {
			if record, ok:=item.(map[string]interface{}); ok {
				if v, ok:=record[field]; ok && v==entityID {
					records[indx]=entity
					couldUpdate=true
					
					fmt.Printf("Updating %s with ID %s \n", entName, entityID)
				}
			}
		}
	}
	if(couldUpdate) {
		err=d.writeAll(records)
	} else {
		return fmt.Errorf("Failed to update. Unable to find any %s record with ID %s", entName, entityID)
	}

	return
}

//Delete the json data based on the id field/value pair
func (d *Driver) Delete(entity Entity) (err error) {
	d.entityDealingWith=entity
	field, entityID:=entity.ID()
	records:= d.Open(entity).Get()
	entName,_:=d.getEntityName()

	couldDelete:=false
	newRecordArray:=make([]interface{},0,0)

	if(len(records)>0){
		for indx,item:= range records {
			if record, ok:=item.(map[string]interface{}); ok {
				if v, ok:=record[field]; ok && v!=entityID {
					records[indx]=entity
					newRecordArray=append(newRecordArray, record)
				} else {
					fmt.Printf("Deleting %s with ID %s \n", entName, entityID)
					couldDelete=true
				}
			}
		}
	}
	if(couldDelete) {
		err=d.writeAll(newRecordArray)
	} else {
		return fmt.Errorf("Failed to delete. Unable to find any %s record with ID %s", entName, entityID)
	}
	return
}