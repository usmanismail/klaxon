#!/bin/sh

files=`gofmt -l ./`
for file in $files
do
	gofmt $file > ${file}_
	mv ${file}_ $file
done