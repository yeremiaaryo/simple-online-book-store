package books

var (
	queryGetBooks = `SELECT id, title, author, isbn, published_date, price
        FROM books`
)
