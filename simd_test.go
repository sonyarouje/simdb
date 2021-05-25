package simdb

import (
	"os"
	"testing"
)

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

func (c Customer) ID() (jsonField string, value interface{}) {
	value = c.CustID
	jsonField = "custid"
	return
}

type Product struct {
	ProdID string  `json:"productId"`
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
}

func (p Product) ID() (jsonField string, value interface{}) {
	value = p.ProdID
	jsonField = "productId"
	return
}

func TestNew(t *testing.T) {
	_, err := New("test")
	if err != nil {
		t.Error(err)
	}
}

func TestInsertCustomer(t *testing.T) {
	driver, err := New("test")
	if err != nil {
		t.Error(err)
	}
	customer := Customer{
		CustID:  "CU1",
		Name:    "sarouje",
		Address: "address",
		Contact: Contact{
			Phone: "45533355",
			Email: "someone@gmail.com",
		},
	}
	err = driver.Insert(customer)
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("./test/Customer"); err != nil {
		t.Errorf("Failed to create customer db file")
	}
	fetched, err := getCustomer(customer)
	if err != nil {
		t.Error(err)
	}
	if fetched.CustID != "CU1" {
		t.Errorf("unable to get the customer with customer id CU1")
	}
}

func TestUpdateCustomer(t *testing.T) {
	driver, err := New("test")
	if err != nil {
		t.Error(err)
	}

	customer := Customer{
		CustID: "CU1",
	}
	if fetched, err := getCustomer(customer); err == nil {
		fetched.Name = "Sony Arouje"
		driver.Update(fetched)

		afterUpdate, _ := getCustomer(customer)
		if afterUpdate.Name != fetched.Name {
			t.Errorf("unable to update the customer record")
		}
	} else {
		t.Error(err)
	}

}

func getCustomer(c Customer) (Customer, error) {
	driver, err := New("test")
	var fetchedCustomer Customer
	err = driver.Open(Customer{}).Where("custid", "=", c.CustID).First().AsEntity(&fetchedCustomer)
	return fetchedCustomer, err
}

func TestInsertProduct(t *testing.T) {
	driver, err := New("test")
	if err != nil {
		t.Error(err)
	}

	product := Product{
		ProdID: "P1",
		Name:   "Test product 1",
		Price:  10,
	}

	err = driver.Insert(product)
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("./test/Product"); err != nil {
		t.Errorf("Failed to create product db file")
	}
	fetched, err := getProduct(product)
	if err != nil {
		t.Error(err)
	}
	if fetched.ProdID != product.ProdID {
		t.Errorf("unable to get the Product with given product id")
	}
	if fetched.Price != product.Price {
		t.Errorf("incorrect price for the fetched product")
	}
}

func getProduct(p Product) (Product, error) {
	driver, err := New("test")
	var fetchedProduct Product
	err = driver.Open(Product{}).Where("productId", "=", p.ProdID).First().AsEntity(&fetchedProduct)
	return fetchedProduct, err
}

func TestDeleteCustomer(t *testing.T) {
	driver, _ := New("test")

	customer := Customer{
		CustID: "CU1",
	}
	if err := driver.Delete(customer); err != nil {
		t.Error(err)
	}
}

func TestDeleteProduct(t *testing.T) {
	driver, _ := New("test")
	product := Product{
		ProdID: "P1",
	}

	if err := driver.Delete(product); err != nil {
		t.Error(err)
	}

}
