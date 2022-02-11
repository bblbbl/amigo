# Amigo

----

#### Software for sql migrations of your project

## Supported databases

-------

* mysql
* postgres

## Installation

-------

```
1. git clone git@github.com:bblbbl/amigo.git

2. cd amigo/bin/amigo

3. go install
``` 

## Usage

-------

#### Create migration

```shell
amigo -create -dir=migration
```

#### Apply migration

```shell
amigo -migrate -dbName=amigo -dbUser=user -dbPassword=password -step=1 -dir=migration
```

#### Rollback migration

```shell
amigo -rollback -dbName=amigo -dbUser=user -dbPassword=password -step=1 -dir=migration
```

#### Single sql script example

```sql
CREATE TABLE user (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(30) NOT NULL
) # dont type ;
```

#### Multi sql script example

```sql
CREATE TABLE user (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(30) NOT NULL
);

CREATE TABLE role (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    role VARCHAR(30) NOT NULL
) # dont type ;
```

## Recommendation

----

In order not to specify keys with access to the database, specify them in the .env file, or use Makefile

## Command Line Arguments

-------

- -dir=some/path Path to the directory with migrations, be default "migration"
- -create Will create two sql files up and down to apply and rollback the migration
- -migrate Applies unapplied migrations
- -rollback Rollback the last migration
- -step=1 Number of migrations to be executed or rolled back
- -dbName=name Database name, required if not specified in .env file
- -dbUser=user Database user, required if not specified in .env file
- -dbPassword=password Database password, required if not specified in .env file
- -dbProvider=postgres Database provider, by default mysql

## Env variables example

-------


```dotenv
DB_NAME=amigo
DB_USER=user
DB_PASSWORD=password
DB_PROVIDER=mysql
```