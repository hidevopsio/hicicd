#!/usr/bin/env bash

openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub