package repository

import (
	"backend/internal/model"
	"context"
	//"database/sql"
	"fmt"
	//"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db DBTX
}

func NewOrderRepository(db DBTX) *OrderRepository {
	return &OrderRepository{db: db}
}

// 注文を作成し、生成された注文IDを返す
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) (string, error) {
	// query := `INSERT INTO orders (user_id, product_id, shipped_status, created_at) VALUES (?, ?, 'shipping', NOW())`
	// result, err := r.db.ExecContext(ctx, query, order.UserID, order.ProductID)
	// if err != nil {
	// 	return "", err
	// }
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return "", err
	// }
	// return fmt.Sprintf("%d", id), nil
	query := `
        INSERT INTO orders (user_id, product_id, shipped_status, created_at)
        VALUES (?, ?, 'shipping', NOW())
    `
    result, err := r.db.ExecContext(ctx, query, order.UserID, order.ProductID)
    if err != nil {
        return "", err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%d", id), nil
}

// 複数の注文IDのステータスを一括で更新
// 主に配送ロボットが注文を引き受けた際に一括更新をするために使用
func (r *OrderRepository) UpdateStatuses(ctx context.Context, orderIDs []int64, newStatus string) error {
	if len(orderIDs) == 0 {
        return nil
    }
    query, args, err := sqlx.In("UPDATE orders SET shipped_status = ? WHERE order_id IN (?)", newStatus, orderIDs)
    if err != nil {
        return err
    }
    query = r.db.Rebind(query)
    _, err = r.db.ExecContext(ctx, query, args...)
    return err

    // query = r.db.Rebind(query)
    // _, err = r.db.ExecContext(ctx, query, args...)
    // return err
	
	// if len(orderIDs) == 0 {
	// 	return nil
	// }
	// query, args, err := sqlx.In("UPDATE orders SET shipped_status = ? WHERE order_id IN (?)", newStatus, orderIDs)
	// if err != nil {
	// 	return err
	// }
	// query = r.db.Rebind(query)
	// _, err = r.db.ExecContext(ctx, query, args...)
	// return err
}

// 配送中(shipped_status:shipping)の注文一覧を取得
func (r *OrderRepository) GetShippingOrders(ctx context.Context) ([]model.Order, error) {


	var orders []model.Order
    query := `
        SELECT
            o.order_id,
            o.product_id,
            p.weight,
            p.value
        FROM orders o
        JOIN products p ON o.product_id = p.product_id
        WHERE o.shipped_status = 'shipping'
        ORDER BY o.created_at ASC
    `
    err := r.db.SelectContext(ctx, &orders, query)
    return orders, err

	// var orders []model.Order
	// query := `
    //     SELECT
    //         o.order_id,
    //         p.weight,
    //         p.value
    //     FROM orders o
    //     JOIN products p ON o.product_id = p.product_id
    //     WHERE o.shipped_status = 'shipping'
    // `
	// err := r.db.SelectContext(ctx, &orders, query)
	// return orders, err
}

// 注文履歴一覧を取得
func (r *OrderRepository) ListOrders(ctx context.Context, userID int, req model.ListRequest) ([]model.Order, int, error) {
 query := `
        SELECT
            o.order_id,
            o.product_id,
            o.shipped_status,
            o.created_at,
            o.arrived_at,
            p.name AS product_name
        FROM orders o
        JOIN products p ON o.product_id = p.product_id
        WHERE o.user_id = ?
    `
    args := []interface{}{userID}

    // 搜索条件
    if req.Search != "" {
        if req.Type == "prefix" {
            query += " AND p.name LIKE ?"
            args = append(args, req.Search+"%")
        } else {
            query += " AND p.name LIKE ?"
            args = append(args, "%"+req.Search+"%")
        }
    }

    // 排序字段
    sortField := "o.order_id"
    switch req.SortField {
    case "product_name":
        sortField = "p.name" // 可以改成 idx_name_id 的索引列
    case "created_at":
        sortField = "o.created_at"
    case "shipped_status":
        sortField = "o.shipped_status"
    case "arrived_at":
        sortField = "o.arrived_at"
    }

    sortOrder := "ASC"
    if strings.ToUpper(req.SortOrder) == "DESC" {
        sortOrder = "DESC"
    }

    query += fmt.Sprintf(" ORDER BY %s %s, o.order_id ASC LIMIT ? OFFSET ?", sortField, sortOrder)
    args = append(args, req.PageSize, req.Offset)

    var orders []model.Order
    if err := r.db.SelectContext(ctx, &orders, query, args...); err != nil {
        return nil, 0, err
    }

    // COUNT 查询（不需要 JOIN）
    countQuery := "SELECT COUNT(*) FROM orders WHERE user_id = ?"
    countArgs := []interface{}{userID}
    if req.Search != "" {
        countQuery += " AND product_id IN (SELECT product_id FROM products WHERE name LIKE ?)"
        if req.Type == "prefix" {
            countArgs = append(countArgs, req.Search+"%")
        } else {
            countArgs = append(countArgs, "%"+req.Search+"%")
        }
    }

    var total int
    if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}
