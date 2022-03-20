#! /bin/bash

DB_NAME=$1

mongo localhost:27017/$DB_NAME validations/db.js