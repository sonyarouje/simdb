package main

import (
	"simd/simd"
	"fmt"
)

type Customer struct {
	Name string `json:"name"`
	Address string `json:"address"`
}

func main(){
	driver, err:=simd.New("mydir")
	if(err!=nil){
		panic(err)
	}

	customer:=Customer {
		Name:"sarouje",
		Address: "address",
	}
	// err=driver.Insert(customer)
	// if(err!=nil){
	// 	panic(err)
	// }

	result:=driver.Open(customer).Where("name","=","sarouje1").Get()
	fmt.Printf("%v", result)
	fmt.Printf("%v", driver.Errors())
	
	// var customers []Customer

	// driver.Write(customers)
}