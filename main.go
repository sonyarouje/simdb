package main

import (
	"simd/db"
	"fmt"
)

type Customer struct {
	Name string `json:"name"`
	Address string `json:"address"`
}

func (c Customer) ID() (jsonField string, value interface{}) {
	value=c.Name
	jsonField="name"
	return
}

func main(){
	driver, err:=db.New("mydir")
	if(err!=nil){
		panic(err)
	}

	// customer:=Customer {
	// 	Name:"sarouje",
	// 	Address: "address",
	// }
	// err=driver.Insert(customer)
	// if(err!=nil){
	// 	panic(err)
	// }

	// result:=driver.Open(customer).Where("name","=","sarouje1").Get()
	// fmt.Printf("%v", result)
	fmt.Printf("%v", driver.Errors())
	
	custUpd:=Customer {
		Name:"sarouje1",
		Address:"UMC2",
	}
	err=driver.Update(custUpd)
	if(err!=nil){
		panic(err)
	}

	custDel:=Customer {
		Name:"sarouje2",
	}
	err=driver.Delete(custDel)
	if(err!=nil){
		panic(err)
	}
	// var customers []Customer

	// driver.Write(customers)
}