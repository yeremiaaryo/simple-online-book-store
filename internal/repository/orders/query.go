package orders

var (
	insertOrderQuery = `
        INSERT INTO orders (user_id, total_amount, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id;
    `
	insertOrderItemQuery = `
        INSERT INTO order_items (order_id, book_id, quantity, price, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?);
    `

	getOrderHistoryByUserID = `
		SELECT id, total_amount, status, created_at, updated_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	getItemsQuery = `
		SELECT id, order_id, book_id, quantity, price
		FROM order_items
		WHERE order_id = ANY(?)
	`
)
