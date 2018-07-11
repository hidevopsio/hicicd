#!/usr/bin/env bash

openssl genrsa -out ssl/app.rsa 1024
openssl rsa -in ssl/app.rsa -pubout > ssl/app.rsa.pub