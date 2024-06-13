package users

var (
	getUsersQuery = `SELECT 
							id, email, password, created_at, updated_at 
						FROM 
						    users`

	insertUserQuery = `INSERT INTO users
							(email, password, created_at, updated_at)
							VALUES(?, ?, ?, ?) RETURNING id;`
)
