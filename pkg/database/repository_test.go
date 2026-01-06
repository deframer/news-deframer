package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetTime(t *testing.T) {
	// 1. Create sqlmock
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	// 2. Connect GORM to the mock
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	// 3. Initialize Repository
	repo := NewFromDB(db)

	// 4. Define Expectations
	now := time.Now()
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow(now))

	// 5. Execute and Assert
	got, err := repo.GetTime()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("expected %v, got %v", now, got)
	}

	// 6. Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
