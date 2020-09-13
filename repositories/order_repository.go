package repositories

import (
	"database/sql"
	"go-product/common"
	"go-product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := "INSERT " + o.table + " SET userID=?, productID=?, orderStatus=?"
	stmt, stmtError := o.mysqlConn.Prepare(sql)
	if stmtError != nil {
		return productID, stmtError
	}

	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderID int64) (isOK bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "DELETE FROM " + o.table + " WHERE ID=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}

	_, err := stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Conn(); err != nil {
		return errConn
	}
	orderID := strconv.FormatInt(order.ID, 10)
	sql := "UPDATE " + o.table + " SET userID=?,productID=?,orderStatus=? WHERE ID=" + orderID
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}

	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)

	return
}

func (o *OrderManagerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}

	sql := "SELECT * FROM " + o.table + " WHERE ID=" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	sql := "SELECT * FROM " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return orderArray, nil
}

func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	sql := "SELECT o.ID,p.productName,o.orderStatus FROM go_product.order AS o LEFT JOIN product as p ON o.productID=p.ID"
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err

}
