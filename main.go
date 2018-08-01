package main

import (
	"simd/db"
	"fmt"
)

type Customer struct {
	Name string `json:"name"`
	Address string `json:"address"`
}

func (c Customer) ID() (value interface{}, jsonField string) {
	value=c.Name
	jsonField="name"
	return
}

func main(){
	driver, err:=db.New("mydir")
	if(err!=nil){
		panic(err)
	}

	customer:=Customer {
		Name:"sarouje",
		Address: "address",
	}
	err=driver.Insert(customer)
	if(err!=nil){
		panic(err)
	}

	result:=driver.Open(customer).Where("name","=","sarouje").Get()
	fmt.Printf("%v", result)
	fmt.Printf("%v", driver.Errors())
	
	// var customers []Customer

	// driver.Write(customers)
}