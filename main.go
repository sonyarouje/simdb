package main

import (
	"simdb/db"
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
	fmt.Println("starting....")

	driver, err:=db.New("mydir")

	if(err!=nil){
		panic(err)
	}

	customer:=Customer {
		CustID:"CUST1",
		Name:"sarouje",
		Address: "address",
		Contact: Contact {
			Phone:"45533355",
			Email:"someone@gmail.com",
		},
	}
	//creates a new Customer file inside the directory passed as the parameter to New()
	//if the Customer file already exist then insert operation will add the customer data to the array
	err=driver.Insert(customer)
	if(err!=nil){
		panic(err)
	}
	
	//GET ALL Customer
	//opens the customer json file and filter all the customers with name sarouje.
	//AsEntity takes an address to Customer array and fills the result to it.
	//we can loop through the customers array and retireve the data.
	var customers []Customer
	err=driver.Open(Customer{}).Where("name","=","sarouje").Get().AsEntity(&customers)
	if(err!=nil){
		panic(err)
	}
	// fmt.Printf("%#v \n", customers)
	
	//GET ONE Customer
	//First() will return the first record from the results 
	//AsEntity takes the address to Customer variable (not an array pointer)
	var customerFirst Customer
	err=driver.Open(Customer{}).Where("custid","=","CUST1").First().AsEntity(&customerFirst)
	if(err!=nil){
		panic(err)
	}
	
	//Update function uses the ID() to get the Id field/value to find the record and update the data.
	customerFirst.Name="Sony Arouje"
	err=driver.Update(customerFirst)
	if(err!=nil){
		panic(err)
	}
	driver.Open(Customer{}).Where("custid","=","CUST1").First().AsEntity(&customerFirst)
	fmt.Printf("%#v \n", customerFirst)

	// Delete
	toDel:=Customer{
		CustID:"CUST1",
	}
	err=driver.Delete(toDel)
	if(err!=nil){
		panic(err)
	}
}