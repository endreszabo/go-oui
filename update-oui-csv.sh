#!/bin/bash
DIR=$(dirname "$0")
rm ${DIR}/oui.csv
wget http://standards-oui.ieee.org/oui/oui.csv -O ${DIR}/oui.csv
