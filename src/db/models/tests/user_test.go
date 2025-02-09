package models_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserService_Create(t *testing.T) {
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

	user := &models.User{
		Name:     new(string),
		Email:    new(string),
		Password: new(string),
        Role:     new(uint),
	}
	*user.Name = "John Doe"
	*user.Email = "john@doe.com"
	*user.Password = "MyPasswd"
    *user.Role = 1

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "users" .* RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()

	userService := &models.UserService{
		Service: mockService,
	}

	if err := userService.Create(user); err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet SQL mock expectations: %v", err)
	}
}

func TestUserService_Fetch(t *testing.T) {
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

	userID := uint(1)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
			AddRow(userID, "John Doe", "john@doe.com", "MyPasswd", 1))

	userService := &models.UserService{
		Service: mockService,
	}

	user, err := userService.Fetch(userID)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	if *user.Name != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got: %s", *user.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet SQL mock expectations: %v", err)
	}
}

func TestUserService_Update(t *testing.T) {
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

    userID := uint(1)
	updatedUser := &models.User{
		Name:     new(string),
		Email:    new(string),
		Password: new(string),
        Role:     new(uint),
	}

    *updatedUser.Name = "Kkk Elba"
    *updatedUser.Email = "kkk@elba.com"
    *updatedUser.Password = "ElbaPass"
    *updatedUser.Role = 1

    mock.ExpectBegin()

    mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
        WithArgs(userID, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
        AddRow(userID, "John Doe", "john@doe.com", "mypass", 0))

    mock.ExpectExec(`UPDATE "users" SET .* WHERE "id" = \$[0-9]+`).
        WithArgs(
            *updatedUser.Name,
            *updatedUser.Email,
            *updatedUser.Password,
            *updatedUser.Role,
            userID,
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))

    mock.ExpectCommit()

    userService := &models.UserService{
        Service: mockService,
    }

    if err := userService.Update(userID, updatedUser); err != nil {
        t.Errorf("Expected no errors, got: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}

func TestUserService_Delete(t *testing.T) {
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

    userID := uint(1)

    mock.ExpectBegin()

    mock.ExpectExec(`DELETE FROM "users" WHERE "users"\."id" = \$1`).
        WithArgs(userID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    mock.ExpectCommit()

    userService := &models.UserService{
        Service: mockService,
    }
    
    if err := userService.Delete(userID); err != nil {
        t.Error("Expected no errors, got: %w", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet SQL mock expectations: %v", err)
    }
}
