#!/bin/sh
printf "$*"
read k
exec "$@"
