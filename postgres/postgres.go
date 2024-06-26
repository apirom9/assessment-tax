package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func NewPostgres(dbUrl string) (*Postgres, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		for i := 1; i <= 5; i++ {
			err = db.Ping()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Postgres{Db: db}, nil
}

func (p *Postgres) UpdateDefaultPersonalDeduction(value float64) error {
	sqlStr := "UPDATE allowance SET allowance_amount=$1 WHERE allowance_type='personal_default'"
	_, err := p.Db.Query(sqlStr, value)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetDefaultPersonalDeduction() (float64, error) {
	result := 0.0
	sqlStr := "SELECT allowance_amount FROM allowance WHERE allowance_type='personal_default'"
	rows, err := p.Db.Query(sqlStr)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (p *Postgres) UpdateMaxKReceipt(value float64) error {
	sqlStr := "UPDATE allowance SET allowance_amount=$1 WHERE allowance_type='kreceipt_max'"
	_, err := p.Db.Query(sqlStr, value)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetMaxKReceipt() (float64, error) {
	result := 0.0
	sqlStr := "SELECT allowance_amount FROM allowance WHERE allowance_type='kreceipt_max'"
	rows, err := p.Db.Query(sqlStr)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
