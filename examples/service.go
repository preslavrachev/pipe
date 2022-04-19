package examples

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/preslavrachev/pipe"
)

type Customer struct {
	/* A Customer model */
}
type Product struct {
	/* A Product model */
}
type Order struct {
	Customer Customer
	Product  Product
}

type Service struct {
	db         *sql.DB
	httpClient *http.Client
}

type sendOrderParams struct {
	customerID string
	productID  string
	customer   Customer
	product    Product
	order      Order
}

func (s *Service) DoSomeComplexThing() error {
	_, err := pipe.New[sendOrderParams]().
		Next(s.loadCustomer).
		Next(s.loadProduct).
		Next(s.sendOrder).
		Do()

	return err
}

func (s *Service) loadCustomer(params *sendOrderParams) error {
	rows, err := s.db.Query("Load the customer", params.customerID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// scan the customer here, potentially returning an error
		rows.Scan(&params.customer)
	}

	return nil
}

func (s *Service) loadProduct(params *sendOrderParams) error {
	rows, err := s.db.Query("Load the product", params.productID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// scan the product here, potentially returning an error
		rows.Scan(&params.product)
	}

	return nil
}

func (s *Service) sendOrder(params *sendOrderParams) error {
	order := Order{
		Customer: params.customer,
		Product:  params.product,
	}

	requestBody, err := json.Marshal(&order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Assuming that everything went successfully
	params.order = order
	return nil
}
