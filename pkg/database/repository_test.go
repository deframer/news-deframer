package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetTime(t *testing.T) {
	// 1. Create sqlmock
	conn, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       conn,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFromDB(db)

	now := time.Now()
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow(now))
	mock.ExpectClose()

	got, err := repo.GetTime()
	assert.NoError(t, err)
	assert.True(t, got.Equal(now))

	assert.NoError(t, conn.Close())

	assert.NoError(t, mock.ExpectationsWereMet())
}
