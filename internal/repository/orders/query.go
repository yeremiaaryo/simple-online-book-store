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
)
