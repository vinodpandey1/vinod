
# Book Rest APIs with Gorilla Mux and MySQL



### Goals

* Create a RESTful web service for a Library. The service must have the following API endpoints:
* (C)reate a new Book
* (R)ead existing Books
* (U)pdate an existing Book
* (D)elete an existing Book


### Prerequisites

* A Book entity has the following properties:
* Author (mandatory)
* Title (mandatory)
* ISBN (mandatory)
* Release Date


* You must have a working Go and MySQL environments.

* Basic familiarity with Go and MySQL.

### About the Application

The application is a simple REST API server that will provide endpoints to allow accessing and manipulating ‘books’.

### API Specification

* Create a new book in response to a valid POST request at /book,

* Update a book in response to a valid PUT request at /book/{id},

* Delete a book in response to a valid DELETE request at /book/{id},

* Fetch a book in response to a valid GET request at /book/{id}, and

* Fetch a list of books in response to a valid GET request at /books.

The {id} will determine which book the request will work with.

### Creating the Database

As our application is simple, we will create only one table called books with the following fields:


Let’s use the following statement to create the database and the table.

    CREATE DATABASE bookdb;
    USE bookdb;
    CREATE TABLE IF NOT EXISTS books
    (
        id INT AUTO_INCREMENT PRIMARY KEY,
        author VARCHAR(50) NOT NULL,
        title VARCHAR(50) NOT NULL,
        isbn VARCHAR(50) NOT NULL,
        release_date DATE DEFAULT (CURRENT_DATE)
    );


### Getting Dependencies

Before we start writing our application, we need to get some dependencies that we will use. We need to get the two following packages:

* mux : The Gorilla Mux router.

* mysql : The MySQL driver.

You can easily use go get to get it:

    go get github.com/gorilla/mux
    go get github.com/go-sql-driver/mysql


### Running Tests
    go test -v

Executing this command should result something like this:

    testing: warning: no tests to run
    PASS
    ok      _/home/book/app 0.051s

### Writing API Tests

