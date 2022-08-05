## ORDER STORAGE SERVICE

This is a simple service which stores order data in database and in-memory cache. NATS streaming is used to improve its performance and reliability.

#### PREREQUISIT

- Postgres database is installed and running on your system
- Configured postgres user and database (see configs/apiserver.toml) 


#### USAGE

  
```
# launch STAN server:
make stan
```
```
# launch service:
make su
```
```
# publish data automatically:
make pu
```
