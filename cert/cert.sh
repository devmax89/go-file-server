#!/bin/bash

# Directory dove generare i certificati
DIR="$(dirname "$0")"

# File di configurazione per le informazioni del certificato
CONFIG="$DIR/cert.conf"

# Crea la directory se non esiste
mkdir -p $DIR

# Genera la chiave privata e il certificato
openssl req -x509 -newkey rsa:4096 -keyout $DIR/key.pem -out $DIR/cert.pem -days 365 -nodes -config $CONFIG