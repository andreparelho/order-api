package order_repository

import "database/sql"

type OrderRepository interface {
	InsertOrder() error
}

type order struct {
	database *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &order{
		database: db,
	}
}

type OrderEntity struct {
}

func (o *order) InsertOrder() error {
	insert, err := o.database.Query("INSERT INTO orders VALUES('23')")
	if err != nil {
		return err
	}
	defer insert.Close()

	return nil
}
