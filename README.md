# goChat

## building from source

You will need a go compiler. Install using your operating systems package manager or download from go.dev and follow the instructions on the website.<br>

You will also need node and make.

To compile just run make<br>
```sh
make
```

### .env file with the following variables is required

- PASSWD_SALT
- SESSION_SECRET
#### these are only needed if using MySql
- MYSQL_USERNAME
- MYSQL_PASSWD
- MYSQL_ADDR

### A mySql database is used by default

Install mySql however you want. Use the provided sql file to set up the database.

### Sqlite

If you want to use sqlite instead of mySql pass the -sqlite flag when running the server
```sh
./goChat -sqlite
```

<br>
<br>
<em>enjoy the spaghetti</em>
