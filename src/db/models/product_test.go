package models

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestProductService_Create(t *testing.T) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create SQL mock: %v", err)
    }
    defer sqlDB.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        PrepareStmt: false,
    })
    if err != nil {
        t.Fatalf("Failed to create GORM DB from SQL mock: %v", err)
    }

    mockService := config.InitMockService(gormDB)

    product := &Product{}
    *product.Name = "Laptop"
    *product.Price = 1200.00
    *product.Points = 120
    *product.Description = "A nice laptop"
    *product.Category = "Electronics"
    // Expect BEGIN transaction
    mock.ExpectBegin()

    // Use ExpectQuery instead of ExpectExec for RETURNING "id"
    mock.ExpectQuery(`INSERT INTO "products" .* RETURNING "id"`).
        WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    // Expect COMMIT transaction
    mock.ExpectCommit()

    productService := &ProductService{
        Service: mockService,
    }

    if err := productService.Create(product); err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}

func TestProductService_Fetch(t *testing.T) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create SQL mock: %v", err)
    }
    defer sqlDB.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        PrepareStmt: false,
    })
    if err != nil {
        t.Fatalf("Failed to create GORM DB from SQL mock: %v", err)
    }

    mockService := config.InitMockService(gormDB)

    productID := uint(1)
    // Correctly match the query and provide 2 arguments (productID and LIMIT)
    mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"\."id" = \$1 ORDER BY "products"\."id" LIMIT \$2`).
        WithArgs(productID, sqlmock.AnyArg()). // Accept both productID and the LIMIT argument
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "points", "category"}).
            AddRow(productID, "Laptop", "A powerful laptop", 1200.50, 100, "Electronics"))

    productService := &ProductService{
        Service: mockService,
    }

    product, err := productService.Fetch(productID)
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }

    if *product.Name != "Laptop" {
        t.Errorf("Expected product name 'Laptop', got: %s", *product.Name)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}


func TestProductService_Update(t *testing.T) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create SQL mock: %v", err)
    }
    defer sqlDB.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        PrepareStmt: false,
    })
    if err != nil {
        t.Fatalf("Failed to create GORM DB from SQL mock: %v", err)
    }

    mockService := config.InitMockService(gormDB)

    productID := uint(1)
    updatedProduct := &Product{}
    *updatedProduct.Name = "Updated laptop"
    *updatedProduct.Description = "new updated laptop"
    *updatedProduct.Price = 1200
    *updatedProduct.Points = 300
    *updatedProduct.Category = "Electronics"

    // Expect BEGIN transaction
    mock.ExpectBegin()

    // Expect SELECT to fetch the existing product
    mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"\."id" = \$1 ORDER BY "products"\."id" LIMIT \$2`).
        WithArgs(productID, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "points", "category"}).
            AddRow(productID, "Laptop", "A powerful laptop", 1200.50, 100, "Electronics"))

    // Fix: Use a more flexible regular expression to match the UPDATE query with dynamic bindings
    mock.ExpectExec(`UPDATE "products" SET .* WHERE "id" = \$[0-9]+`).
        WithArgs(
            updatedProduct.Name,
            updatedProduct.Description,
            updatedProduct.Price,
            updatedProduct.Points,
            updatedProduct.Category,
            productID,
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Expect COMMIT transaction
    mock.ExpectCommit()

    productService := &ProductService{
        Service: mockService,
    }

    if err := productService.Update(productID, updatedProduct); err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}

func TestProductService_Delete(t *testing.T) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create SQL mock: %v", err)
    }
    defer sqlDB.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        PrepareStmt: false,
    })
    if err != nil {
        t.Fatalf("Failed to create GORM DB from SQL mock: %v", err)
    }

    mockService := config.InitMockService(gormDB)

    productID := uint(1)

    // Expect BEGIN transaction
    mock.ExpectBegin()

    // Expect DELETE query
    mock.ExpectExec(`DELETE FROM "products" WHERE "products"\."id" = \$1`).
        WithArgs(productID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Expect COMMIT transaction
    mock.ExpectCommit()

    productService := &ProductService{
        Service: mockService,
    }

    if err := productService.Delete(productID); err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}

