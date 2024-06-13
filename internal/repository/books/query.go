package books

var (
	queryGetBooks = `SELECT id, title, author, isbn, published_date, price, created_at, updated_at
        FROM books`
)
