package main // Package main, Do not change this line.

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Product represents the structure for a product entity
type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// SupplyChainContract defines the smart contract structure
type SupplyChainContract struct {
	contractapi.Contract
}

// getTimestamp returns the transaction timestamp as a string
func (s *SupplyChainContract) getTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	return time.Unix(txTimestamp.Seconds, int64(txTimestamp.Nanos)).Format(time.RFC3339), nil
}

// InitLedger initializes the ledger with some example products
func (s *SupplyChainContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	timestamp, err := s.getTimestamp(ctx)
	if err != nil {
		return err
	}

	// Initial set of products to populate the ledger
	products := []Product{
		{ID: "p1", Name: "Laptop", Status: "Manufactured", Owner: "CompanyA", CreatedAt: timestamp, UpdatedAt: timestamp, Description: "High-end gaming laptop", Category: "Electronics"},
		{ID: "p2", Name: "Smartphone", Status: "Manufactured", Owner: "CompanyB", CreatedAt: timestamp, UpdatedAt: timestamp, Description: "Latest model smartphone", Category: "Electronics"},
	}

	for _, product := range products {
		if err := s.putProduct(ctx, &product); err != nil {
			return err
		}
		//log.Printf("id: %s name: %s status: %s owner: %s createdAt: %s updatedAt %s description: %s category: %s", product.ID, product.Name, product.Status, product.Owner, product.CreatedAt, product.UpdatedAt, product.Description, product.Category)
	}

	return nil
}

// CreateProduct creates a new product in the ledger
func (s *SupplyChainContract) CreateProduct(ctx contractapi.TransactionContextInterface, id, name, owner, description, category string) error {
	// Write your implementation here
	exists, err := s.ProductExists(ctx, id)
	//log.Printf("exists: %v", exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Product %v already exists", id)
	}

	if id == "" {
		return fmt.Errorf("Product id cannot be empty")
	}
	// if name == "" {
	// 	return fmt.Errorf("Product name cannot be empty")
	// }
	// if owner == "" {
	// 	return fmt.Errorf("Product owner cannot be empty")
	// }
	// if description == "" {
	// 	return fmt.Errorf("Product description cannot be empty")
	// }
	// if category == "" {
	// 	return fmt.Errorf("Product category cannot be empty")
	// }

	timestamp, err := s.getTimestamp(ctx)
	if err != nil {
		return err
	}

	product := &Product{ID: id, Name: name, Status: "Manufactured", Owner: owner, CreatedAt: timestamp, UpdatedAt: timestamp, Description: description, Category: category}

	if err := s.putProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to put product: %v", err)
	}
	//log.Printf("id: %s name: %s status: %s owner: %s createdAt: %s updatedAt %s description: %s category: %s", product.ID, product.Name, product.Status, product.Owner, product.CreatedAt, product.UpdatedAt, product.Description, product.Category)
	exists, err = s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	//log.Printf("exists: %v", exists)
	if !exists {
		return fmt.Errorf("Product creation failed")
	}

	return nil
}

// UpdateProduct allows updating a product's status, owner, description, and category
func (s *SupplyChainContract) UpdateProduct(ctx contractapi.TransactionContextInterface, id string, newStatus string, newOwner string, newDescription string, newCategory string) error {
	// Write your implementation here
	if id == "" {
		return fmt.Errorf("Product id cannot be empty")
	}

	//log.Printf("status: %s owner: %s description: %s category: %s", newStatus, newOwner, newDescription, newCategory)
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		s.CreateProduct(ctx, id, "", "", "", "")
	}
	product, err := s.QueryProduct(ctx, id)
	//log.Printf("error: %v", err)
	if err != nil {
		return err
	}

	if newStatus != "" {
		product.Status = newStatus
	}
	if newOwner != "" {
		product.Owner = newOwner
	}
	if newDescription != "" {
		product.Description = newDescription
	}
	if newCategory != "" {
		product.Category = newCategory
	}

	timestamp, err := s.getTimestamp(ctx)
	if err != nil {
		return err
	}
	product.UpdatedAt = timestamp

	s.putProduct(ctx, product)
	//log.Printf("error: %v", err)
	if err := s.putProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}

	return nil
}

// TransferOwnership changes the owner of a product
func (s *SupplyChainContract) TransferOwnership(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	// Write your implementation here
	if id == "" {
		return fmt.Errorf("Product id cannot be empty")
	}
	if newOwner == "" {
		return fmt.Errorf("New owner cannot be empty")
	}
	exists, err := s.ProductExists(ctx, id)
	if !exists {
		return fmt.Errorf("Product does not exist: %v", err)
	}
	product, err := s.QueryProduct(ctx, id)
	if err != nil {
		return err
	}
	timestamp, err := s.getTimestamp(ctx)
	if err != nil {
		return err
	}

	product.Owner = newOwner
	product.UpdatedAt = timestamp

	s.putProduct(ctx, product)
	return nil
}

// QueryProduct retrieves a single product from the ledger by ID
func (s *SupplyChainContract) QueryProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	// Write your implementation here
	exists, err := s.ProductExists(ctx, id)
	if !exists {
		return nil, fmt.Errorf("Product does not exist: %v", err)
	}

	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	//log.Printf("productJSON: %v", productJSON)
	var product *Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal product: %v", err)
	}
	return product, nil
}

// putProduct is a helper method for inserting or updating a product in the ledger
func (s *SupplyChainContract) putProduct(ctx contractapi.TransactionContextInterface, product *Product) error {
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(product.ID, productJSON)
}

// ProductExists is a helper method to check if a product exists in the ledger
func (s *SupplyChainContract) ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return productJSON != nil, nil
}

// GetAllProducts is a helper method to retrieve all products from the ledger
func (s *SupplyChainContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		if err := json.Unmarshal(queryResponse.Value, &product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SupplyChainContract{})
	if err != nil {
		fmt.Printf("Error creating supply chain chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting supply chain chaincode: %s", err.Error())
	}
}
