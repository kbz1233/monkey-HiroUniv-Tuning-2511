package repository

import (
	"backend/internal/model"
	"context"
	"strings"
	"fmt"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) *ProductRepository {
	return &ProductRepository{db: db}
}

// 商品一覧を全件取得し、アプリケーション側でページング処理を行う
func (r *ProductRepository) ListProducts(ctx context.Context, userID int, req model.ListRequest) ([]model.Product, int, error) {
	 var products []model.Product
    var total int

    baseQuery := "FROM products WHERE 1=1"
    args := []interface{}{}

    if req.Search != "" {
        if req.Type == "prefix" {
            baseQuery += " AND name LIKE ?"
            args = append(args, req.Search+"%")
        } else {
            baseQuery += " AND name LIKE ?"
            args = append(args, "%"+req.Search+"%")
        }
    }

    // 总数
    countQuery := "SELECT COUNT(*) " + baseQuery
    if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
        return nil, 0, err
    }

    // 排序字段
    sortField := "product_id"
    sortOrder := "ASC"
    if req.SortField != "" {
        sortField = req.SortField
    }
    if strings.ToUpper(req.SortOrder) == "DESC" {
        sortOrder = "DESC"
    }

    dataQuery := fmt.Sprintf(`
        SELECT product_id, name, value, weight, image, description
        %s
        ORDER BY %s %s, product_id ASC
        LIMIT ? OFFSET ?
    `, baseQuery, sortField, sortOrder)

    args = append(args, req.PageSize, req.Offset)

    if err := r.db.SelectContext(ctx, &products, dataQuery, args...); err != nil {
        return nil, 0, err
    }

    return products, total, nil


	//  var products []model.Product
    // var total int

    // // 允许排序字段列表，防止 SQL 注入
    // allowedSortFields := map[string]string{
    //     "product_id":  "product_id",
    //     "name":        "name",
    //     "value":       "value",
    //     "weight":      "weight",
    // }

    // sortField := "product_id"
    // if f, ok := allowedSortFields[req.SortField]; ok {
    //     sortField = f
    // }

    // sortOrder := "ASC"
    // if strings.ToUpper(req.SortOrder) == "DESC" {
    //     sortOrder = "DESC"
    // }

    // // 搜索条件
    // baseQuery := "FROM products WHERE 1=1"
    // args := []interface{}{}
    // if req.Search != "" {
    //     baseQuery += " AND (name LIKE ? OR description LIKE ?)"
    //     pattern := "%" + req.Search + "%"
    //     args = append(args, pattern, pattern)
    // }

    // // 总数查询
    // countQuery := "SELECT COUNT(*) " + baseQuery
    // if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
    //     return nil, 0, err
    // }

    // // 数据查询
    // dataQuery := fmt.Sprintf(`
    //     SELECT product_id, name, value, weight, image, description
    //     %s
    //     ORDER BY %s %s, product_id ASC
    //     LIMIT ? OFFSET ?
    // `, baseQuery, sortField, sortOrder)

    // args = append(args, req.PageSize, req.Offset)

    // if err := r.db.SelectContext(ctx, &products, dataQuery, args...); err != nil {
    //     return nil, 0, err
    // }

    // return products, total, nil
	
	// var products []model.Product
	// var total int

	// // ベースクエリ
	// baseQuery := `
	// 	FROM products
	// 	WHERE 1=1
	// `
	// args := []interface{}{}

	// // 検索条件
	// if req.Search != "" {
	// 	baseQuery += " AND (name LIKE ? OR description LIKE ?)"
	// 	searchPattern := "%" + req.Search + "%"
	// 	args = append(args, searchPattern, searchPattern)
	// }

	// // 総件数取得
	// countQuery := "SELECT COUNT(*) " + baseQuery
	// if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
	// 	return nil, 0, err
	// }

	// // 並び順を指定（安全のためデフォルトを設定）
	// sortField := "product_id"
	// sortOrder := "ASC"
	// if req.SortField != "" {
	// 	sortField = req.SortField
	// }
	// if strings.ToUpper(req.SortOrder) == "DESC" {
	// 	sortOrder = "DESC"
	// }

	// // ページング付きデータ取得
	// dataQuery := fmt.Sprintf(`
	// 	SELECT product_id, name, value, weight, image, description
	// 	%s
	// 	ORDER BY %s %s, product_id ASC
	// 	LIMIT ? OFFSET ?
	// `, baseQuery, sortField, sortOrder)

	// args = append(args, req.PageSize, req.Offset)

	// if err := r.db.SelectContext(ctx, &products, dataQuery, args...); err != nil {
	// 	return nil, 0, err
	// }

	// return products, total, nil
	// var products []model.Product
	// baseQuery := `
	// 	SELECT product_id, name, value, weight, image, description
	// 	FROM products
	// `
	// args := []interface{}{}

	// if req.Search != "" {
	// 	baseQuery += " WHERE (name LIKE ? OR description LIKE ?)"
	// 	searchPattern := "%" + req.Search + "%"
	// 	args = append(args, searchPattern, searchPattern)
	// }

	// baseQuery += " ORDER BY " + req.SortField + " " + req.SortOrder + " , product_id ASC"

	// err := r.db.SelectContext(ctx, &products, baseQuery, args...)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// total := len(products)
	// start := req.Offset
	// end := req.Offset + req.PageSize
	// if start > total {
	// 	start = total
	// }
	// if end > total {
	// 	end = total
	// }
	// pagedProducts := products[start:end]

	// return pagedProducts, total, nil
}
