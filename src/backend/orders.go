package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Total        int         `json:"total"`
	Packs        []orderPack `json:"order_packs"`
}

type orderPack struct {
	OrderID  int `json:"order_id"`
	PackID   int `json:"pack_id"`
	Quantity int `json:"quantity"`
}

type pack struct {
	ID   int `json:"id"`
	size int `json:"size"`
}

func getOrders(db *sql.DB) ([]order, error) {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order
	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Total); err != nil {
			return nil, err
		}

		err = o.getOrderPacks(db)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (o *order) createOrder(db *sql.DB) error {

	err := o.calculatePacksRequired(db)

	res, err := db.Exec("INSERT INTO orders(customerName, total) VALUES (?, ?)", o.CustomerName, o.Total)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	o.ID = int(id)

	return nil
}

func (o *order) calculatePacksRequired(db *sql.DB) error {
	var packs, _ = getPacks(db)
	var total = o.Total
	var packAmounts [5]int
	output := make(map[int]int)

	for i := len(packs) - 1; i > 0; i-- {
		for total >= packs[i].size {
			total = total - packs[i].size
			packAmounts[i]++
		}

		if (total > packs[i-1].size) && (packs[i].size-total < packs[0].size) {
			total = total - packs[i].size
			packAmounts[i]++
		}
	}

	if total > 0 {
		for total > 0 {
			total = total - packs[0].size
			packAmounts[0]++
		}
	}

	for i := 0; i < len(packs); i++ {
		if packAmounts[i] != 0 {
			output[packs[i].size] = packAmounts[i]
		}
	}
	fmt.Println(output)
	return nil

}

func (o *order) getOrderPacks(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM order_packs WHERE order_id = ?", o.ID)
	if err != nil {
		return err // just return the error (not nil)
	}
	defer rows.Close()

	var orderPacks []orderPack
	for rows.Next() {
		var op orderPack
		if err := rows.Scan(&op.OrderID, &op.PackID, &op.Quantity); err != nil {
			return err
		}
		orderPacks = append(orderPacks, op)
	}
	o.Packs = orderPacks
	return nil
}

func getPacks(db *sql.DB) ([]pack, error) {
	rows, err := db.Query("SELECT * FROM packs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packs []pack
	for rows.Next() {
		var p pack
		if err := rows.Scan(&p.ID, &p.size); err != nil {
			return nil, err
		}

		packs = append(packs, p)
	}
	return packs, nil
}
