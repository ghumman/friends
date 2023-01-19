#!/bin/sh
./mvnw package;
java -jar target/*.jar;
