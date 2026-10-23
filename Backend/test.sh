#!/bin/bash
export $(cat .env | xargs) && go test -count=1 ./modules/**/application/