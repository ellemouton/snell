# Snell :zap:

https://snell.ellemouton.com/

Using Lightning and LSATs to create a blog website.
- pay per article
- no user registration
- once you have paid for an article, you will have the necessary LSAT and will not need to pay again.

----------------------------------------------------

## Setup:

1. Install:

go get -u github.com/ellemouton/snell

2. Back end services that need to be running:

- lnd
- etcd

3. Prep the DB:

I have used a mysql database for this project

- create a database called 'snell'
- create the tables using the schema in snell/db/schema.sql

4. Start Snell:

Use '$ snell --help' to see the various flags that can be set and the default values that are used if the flags are not set. The dafault values should be fine for local dev and so you can just use "$ snell" to run it. 

An example of how to run it with specific flags:

$ snell  --lnd_cert=/home/admin/.lnd/tls.cert --mac_path=/home/admin/.lnd/data/chain/bitcoin/mainnet/admin.macaroon --db_address=username:password@(127.0.0.1:3306)/snell?
