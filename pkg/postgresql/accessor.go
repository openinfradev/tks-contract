package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/sktelecom/tks-contract/pkg/log"

	_ "github.com/lib/pq"
)

// Accessor is an accessor  for PostgresqlDB.
type Accessor struct {
	db *sql.DB
}

// New returns a new Postgresql.
func New(db *sql.DB) *Accessor {
	return &Accessor{
		db: db,
	}
}

// Close closes database session.
func (p *Accessor) Close() error {
	p.db.Close()
	return nil
}

// Get returns result of querying from DB.
// Support both non-transactional and transactional queries.
func (p *Accessor) Get(tx *sql.Tx, fields, table string, conditions map[string]interface{}) (*sql.Rows, error) {
	var query string
	if len(conditions) != 0 {
		conditionSql := getVarSyntaxFromMaps(conditions)
		query = fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, fields, table, conditionSql[0])
	} else {
		query = fmt.Sprintf(`SELECT %s FROM %s`, fields, table)
	}

	if tx == nil {
		return p.db.Query(query)
	}

	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return rows, nil
}

// Insert inserts new column into table.
func (p *Accessor) Insert(tx *sql.Tx, table string, values ...interface{}) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s VALUES(%s)`, table, getVarSyntax(len(values)))
	var (
		res sql.Result
		err error
	)
	if tx == nil {
		if res, err = p.db.Exec(query, values...); err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	if res, err = tx.Exec(query, values...); err != nil {
		tx.Rollback()
		return 0, err
	}
	return res.RowsAffected()
}

// Delete deletes a row which meets a condition in table.
func (p *Accessor) Delete(tx *sql.Tx, table string, conditions map[string]interface{}) (int64, error) {
	var (
		res sql.Result
		err error
	)
	conditionKeySql := getVarSyntaxFromMaps(conditions)
	conditionValues := getValueSliceFromMaps(conditions)

	query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, table, conditionKeySql[0])
	if tx == nil {
		if res, err = p.db.Exec(query, conditionValues...); err != nil {
			log.Fatal(err)
			return 0, err
		}
		return res.RowsAffected()
	}

	if res, err = tx.Exec(query, conditionValues...); err != nil {
		tx.Rollback()
		return 0, err
	}
	return res.RowsAffected()
}

// Update updates values of specific row which meets a condition in table.
func (p *Accessor) Update(tx *sql.Tx, table string, values, conditions map[string]interface{}) (int64, error) {
	var (
		res sql.Result
		err error
	)
	sqlArr := getVarSyntaxFromMaps(values, conditions)
	args := getValueSliceFromMaps(values, conditions)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, sqlArr[0], sqlArr[1])
	if tx == nil {
		if res, err = p.db.Exec(query, args...); err != nil {
			log.Fatal(err)
			return 0, err
		}
		return res.RowsAffected()
	}
	if res, err = tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return 0, err
	}
	return res.RowsAffected()
}

func (p *Accessor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *Accessor) BeginTx() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *Accessor) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

// getVarSyntax makes "$1, $2, $3..." string for SQL query.
func getVarSyntax(count int) string {
	var (
		result string
		idx    int
		start  bool = true
	)
	for idx = 1; idx <= count; idx++ {
		if !start {
			result += ", "
		}
		result += fmt.Sprintf(`$%d`, idx)
		start = false
	}
	return result
}

// getValuesSliceFromMaps returns one slice gathering all values from multiple maps.
func getValueSliceFromMaps(maps ...map[string]interface{}) []interface{} {
	result := make([]interface{}, 0)
	for i := range maps {
		for _, v := range maps[i] {
			result = append(result, v)
		}
	}
	return result
}

// getVarSyntaxFromMaps returns multiple varSyntax "name=$1, id=$2 ..." from multiple maps.
// Index of varSyntax between multiple maps increases continously.
func getVarSyntaxFromMaps(maps ...map[string]interface{}) []string {
	var (
		idx    int = 1
		result []string
	)
	for i := range maps {
		var (
			temp  string
			start bool = true
		)
		for k := range maps[i] {
			if !start {
				temp += ", "
			}
			temp += fmt.Sprintf("%s=$%d", k, idx)
			start = false
			idx++
		}
		result = append(result, temp)
	}
	return result
}
