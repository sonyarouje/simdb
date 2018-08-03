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
	// tmpCust:=&Customer{}
	err=driver.Open(Customer{}).Where("custid","=","CU1").First().AsEntity(&custOut)
	fmt.Printf("%s %s", custOut.Name, custOut.Contact.Email)

	fmt.Println("")
	var custArray []Customer
	err=driver.Open(customer).Where("name","=","sarouje").Get().AsEntity(&custArray)
	// driver.ToEntityArray(result2, &custArray)
	fmt.Printf("%#v", custArray)

	// var entityArray []Customer
	// var custOut1 Customer;
	// for _, item:= range result2 {
	// 	driver.ToEntity(item, &custOut1)
	// 	entityArray=append(entityArray, custOut1)
	// }
	// fmt.Printf("%#v", entityArray)

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