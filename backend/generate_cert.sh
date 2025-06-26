#!/bin/bash

# Generate self-signed certificate for localhost with SAN extension
# This is required for modern browsers to accept the certificate

# Create a config file for the certificate
cat > cert.conf <<EOF
[req]
default_bits = 4096
prompt = no
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C = US
ST = CA
L = San Francisco
O = Muse Development
OU = Development
CN = 127.0.0.1

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = 127.0.0.1
IP.1 = 127.0.0.1
EOF

# Generate the certificate with the config
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes -config cert.conf -extensions v3_req

# Clean up the config file
rm cert.conf

echo "Generated cert.pem and key.pem for HTTPS development with SAN extension"
echo "Certificate includes localhost and 127.0.0.1 in Subject Alternative Names"
echo "Add these to your .gitignore to keep them local" 