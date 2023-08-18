#!/bin/bash

cd "$(dirname "$0")/server" || exit
./go-url-shortener run
