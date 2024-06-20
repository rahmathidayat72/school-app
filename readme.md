## Notes

- Genered go mod

  ```
  go mod init projectname
  ```

- Genered fremwork Echo

  ```
  go get github.com/labstack/echo/v4
  ```

- Genered Gorm

  ```
  go get -u gorm.io/gorm
  ```

- Genered driver Gorm postgres

  ```
  go get -u gorm.io/driver/postgres
  ```

- Genered Godotenv (to load .evn automatically)

  ```
  go get -u github.com/joho/godotenv
  ```

- Created file `local.env`
  ```
  export DBUSER='DB_username'
  export DBPASS='DB_password'
  export DBHOST='localhost'
  export DBPORT='3306'
  export DBNAME='db_name'
  ```
