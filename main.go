package main

import (
	"simd/db"
	"fmt"
)

type Customer struct {
	CustID string `json:"custid"`
	Name string `json:"name"`
	Address string `json:"address"`
	Contact Contact
}

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (c Customer) ID() (jsonField string, value interface{}) {
	value=c.CustID
	jsonField="custid"
	return
}

func main(){
	driver, err:=db.New("mydir")
	if(err!=nil){
		panic(err)
	}

	customer:=Customer {
		CustID:"CU1",
		Name:"sarouje",
		Address: "address",
		Contact: Contact {
			Phone:"45533355",
			Email:"someone@gmail.com",
		},
	}
	// err=driver.Insert(customer)
	// if(err!=nil){
	// 	panic(err)
	// }

	var custOut Customer;

	// result:=driver.Open(customer).Where("custid","=","CU1").First()
	// tmp:=driver.ToEntity(result, &custOut)
	// fmt.Printf("%#v", tmp)

	// fmt.Printf("%s %s", custOut.Name, custOut.Contact.Email)

	result2:=driver.Open(customer).Where("name","=","sarouje").Get()
	entityArray:=driver.ToEntityArray(result2, &custOut)
	fmt.Printf("%#v", entityArray)

	// custUpd:=Customer {
	// 	Name:"sarouje1",
	// 	Address:"UMC2",
	// }
	// err=driver.Update(custUpd)
	// if(err!=nil){
	// 	panic(err)
	// }

	// custDel:=Customer {
	// 	Name:"sarouje2",
	// }
	// err=driver.Delete(custDel)
	// if(err!=nil){
	// 	panic(err)
	// }
	// var customers []Customer

	// driver.Write(customers)
}