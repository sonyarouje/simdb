# simdb
A simple json db in GO

I sometimes write programs for RaspberryPi using nodejs and use a json file as a data storage. There are so many libraries in nodejs to deal with json file as a data storage. I could'nt find a similar library in GO to use in GO based RPi projects. So decided to write one.

Keep in mind this library can be used in less data intensive applications. I normally use these kind of json data to store execution rules for sensors, etc.

**Go Doc** https://godoc.org/github.com/sonyarouje/simdb/db

**What it does?**
This library enables to store, retrieve, update and delete data from the json db.

**Usage**
Let's have a look at, how to store some data and manipulate them using simd.

```go
package main

import "github.com/sonyarouje/simdb/db"

type Customer struct {
	CustID  string `json:"custid"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact Contact
}

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// ID any struct that needs to persist should implement this function defined
// in Entity interface.
func (c Customer) ID() (jsonField string, value interface{}) {
	value = c.CustID
	jsonField = "custid"
	return
}

func main() {
	driver, err := db.New("data")
	if err != nil {
		panic(err)
	}

	customer := Customer{
		CustID:  "CUST1",
		Name:    "sarouje",
		Address: "address",
		Contact: Contact{
			Phone: "45533355",
			Email: "someone@gmail.com",
		},
	}

	//creates a new Customer file inside the directory passed as the parameter to New()
	//if the Customer file already exist then insert operation will add the customer data to the array
	err = driver.Insert(customer)
	if err != nil {
		panic(err)
	}

	//GET ALL Customer
	//opens the customer json file and filter all the customers with name sarouje.
	//AsEntity takes a pointer to Customer array and fills the result to it.
	//we can loop through the customers array and retireve the data.
	var customers []Customer
	err = driver.Open(Customer{}).Where("name", "=", "sarouje").Get().AsEntity(&customers)
	if err != nil {
		panic(err)
	}

	//GET ONE Customer
	//First() will return the first record from the results
	//AsEntity takes a pointer to Customer variable (not an array pointer)
	var customerFrist Customer
	err = driver.Open(Customer{}).Where("custid", "=", "CUST1").First().AsEntity(&customerFrist)
	if err != nil {
		panic(err)
	}

	//Update function uses the ID() to get the Id field/value to find the record and update the data.
	customerFrist.Name = "Sony Arouje"
	err = driver.Update(customerFrist)
	if err != nil {
		panic(err)
	}

	//Delete
	toDel := Customer{
		CustID: "CUST1",
	}
	err = driver.Delete(toDel)
}
```



Started as a library to learn GO.

Some of the codes to apply json filtering are borrowed from https://github.com/thedevsaddam/gojsonq
