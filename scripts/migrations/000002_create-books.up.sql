CREATE TABLE IF NOT EXISTS books (
    id SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn TEXT NOT NULL UNIQUE,
    published_date DATE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Index for faster searching by author
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);

-- Index for faster searching by title
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);

INSERT INTO books (title, author, isbn, published_date, price, created_at, updated_at) VALUES
    ('The Catcher in the Rye', 'J.D. Salinger', '9780316769488', '1951-07-16', 10.99, 1620993600000, 1620993600000),
    ('To Kill a Mockingbird', 'Harper Lee', '9780061120084', '1960-07-11', 7.99, 1620993600000, 1620993600000),
    ('1984', 'George Orwell', '9780451524935', '1949-06-08', 9.99, 1620993600000, 1620993600000),
    ('Pride and Prejudice', 'Jane Austen', '9780141040349', '1813-01-28', 6.99, 1620993600000, 1620993600000),
    ('The Great Gatsby', 'F. Scott Fitzgerald', '9780743273565', '1925-04-10', 10.99, 1620993600000, 1620993600000),
    ('Moby Dick', 'Herman Melville', '9781503280786', '1851-10-18', 8.99, 1620993600000, 1620993600000),
    ('War and Peace', 'Leo Tolstoy', '9781400079988', '1869-01-01', 12.99, 1620993600000, 1620993600000),
    ('Crime and Punishment', 'Fyodor Dostoevsky', '9780486415871', '1866-01-01', 11.99, 1620993600000, 1620993600000),
    ('The Hobbit', 'J.R.R. Tolkien', '9780547928227', '1937-09-21', 8.99, 1620993600000, 1620993600000),
    ('Catch-22', 'Joseph Heller', '9781451626650', '1961-11-10', 9.99, 1620993600000, 1620993600000);

