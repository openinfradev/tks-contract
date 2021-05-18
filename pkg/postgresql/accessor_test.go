package postgresql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sktelecom/tks-contract/pkg/postgresql"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()
	accessor := postgresql.New(db)
	defer accessor.Close()

	mock.ExpectExec("INSERT INTO records").
		WithArgs("gopher", "828cec77-1da5-4ba6-90d0-270d71be3c55", 50).
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := accessor.Insert(nil, "records(name, id, score)",
		"gopher",
		"828cec77-1da5-4ba6-90d0-270d71be3c55",
		50)

	if err != nil {
		t.Errorf("error was not expected while creating contract: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulilled expectations: %s", err)
	}
	t.Log("updated count: ", count)
}

func TestInsertWithTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()
	accessor := postgresql.New(db)
	defer accessor.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO records").
		WithArgs("gopher", "828cec77-1da5-4ba6-90d0-270d71be3c55", 50).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := accessor.BeginTx()
	count, err := accessor.Insert(tx, "records(name, id, score)",
		"gopher",
		"828cec77-1da5-4ba6-90d0-270d71be3c55",
		50)
	tx.Commit()
	if err != nil {
		t.Errorf("error was not expected while creating contract: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulilled expectations: %s", err)
	}
	t.Log("updated count: ", count)
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"name", "id", "score"}).
		AddRow("gopher", "828cec77-1da5-4ba6-90d0-270d71be3c55", 50)
	mock.ExpectQuery("^SELECT (.+) FROM records$").WillReturnRows(mockRows)

	accessor := postgresql.New(db)
	// have to call Close()
	defer accessor.Close()

	rows, err := accessor.Get(nil, "*", "records", map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, id string
			score    int
		)
		rows.Scan(&name, &id, &score)
		t.Logf("scanned row: %s %s %d", name, id, score)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE records").
		WithArgs(70, "8192039", "gopher").
		WillReturnResult(sqlmock.NewResult(1, 1))
	accessor := postgresql.New(db)
	// have to call Close()
	defer accessor.Close()
	updateValues := map[string]interface{}{
		"score": 70,
	}
	conditionValues := map[string]interface{}{
		"id":   "8192039",
		"name": "gopher",
	}
	count, err := accessor.Update(nil, "records", updateValues, conditionValues)
	if err != nil {
		t.Error(err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	t.Logf("updated count: %d", count)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection.", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM records").
		WithArgs("gopher", 31, 50).
		WillReturnResult(sqlmock.NewResult(1, 1))
	accessor := postgresql.New(db)
	// have to call Close()
	defer accessor.Close()
	condition := map[string]interface{}{
		"name":  "gopher",
		"age":   31,
		"score": 50,
	}
	count, err := accessor.Delete(nil, "records", condition)
	if err != nil {
		t.Error(err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	t.Logf("updated count: %d", count)
}
