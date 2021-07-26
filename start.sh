#!/bin/bash

case $1 in
	test)
		go test -v -cover=true
		;;
	help|*)
		echo "$0 [test|help]"
		;;
esac
