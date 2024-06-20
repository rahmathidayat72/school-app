## Notes

* Genered go mod
    ``````
    go mod init projectname
    ``````

* Genered fremwork Echo
    ``````
    go get github.com/labstack/echo/v4
    ``````

* Genered Gorm
    ``````
    go get -u gorm.io/gorm
    ``````

* Genered driver Gorm mysql
    ``````
    go get -u gorm.io/driver/mysql
    ``````

* Genered viper (to load .evn automatically)
    ``````
    go get -u github.com/spf13/viper
    ``````

* Created file `local.env`
    ``````
    export DBUSER='DB_username'
    export DBPASS='DB_password'
    export DBHOST='localhost'
    export DBPORT='3306'
    export DBNAME='db_name'
    ``````