package models

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func TestProductView_Create(t *testing.T) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create SQL mock: %v", err)
    }
    defer sqlDB.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        // Make sure prepared statements are not enabled so that GORM directly executes the SQL.
        PrepareStmt: false,
    })
    if err != nil {
        t.Fatalf("Failed to create GORM DB from SQL mock: %v", err)
    }

    // Assuming your InitMockService wraps gormDB in your service implementation.
    mockService := config.InitMockService(gormDB)

    product := &Product{
        Name:        "Laptop",
        Description: "A powerful laptop",
        Price:       1200.50,
        Points:      100,
        Category:    "Electronics",
    }

    // Remove ExpectPrepare; only expect Exec
    mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "products"`)).
        WithArgs(product.Name, product.Description, product.Price, product.Points, product.Category).
        WillReturnResult(sqlmock.NewResult(1, 1))

    productView := &ProductView{
        Product: product,
        service: mockService,
    }

    if err := productView.Create(); err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %v", err)
    }
}

