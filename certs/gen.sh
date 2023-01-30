#!/usr/bin/env bash

# CA config
cat << EOF > ca-config.json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "intermediate": {
        "is_ca": true,
        "ca_constraint": {
          "is_ca": true,
          "max_path_len": 1,
          "man_path_len_zero": false
        },
        "expiry": "87600h",
        "usages": [
          "key signature",
          "signing",
          "server auth",
          "client auth"
        ]
      },
      "server": {
        "expiry": "87600h",
        "usages": [
          "signing",
          "server",
          "client"
        ]
      },
      "client": {
        "expiry": "87600h",
        "usages": [
          "signing",
          "client"
        ]
      }
    }
  }
}
EOF

# signed CA request
cat << EOF > ca-csr.json
{
  "CN": "grpc",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "ca": {
    "expiry": "87600h"
  },
  "names": [
    {
      "C":  "CN",
      "L":  "XA",
      "ST": "XA",
      "O":  "grpc",
      "OU": "grpc"
    }
  ]
}
EOF

# signed server request
cat << EOF > server-csr.json
{
  "CN": "grpc.server",
  "hosts": [
    "grpc.server",
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C":  "CN",
      "L":  "XA",
      "ST": "XA",
      "O":  "grpc",
      "OU": "grpc"
    }
  ]
}
EOF

# CA certificate
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
echo "sign CA certificate successfully"

# server certificate
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile server server-csr.json | cfssljson -bare server
echo "sign CA server certificate successfully"

cat ca.pem >> server.pem

# signed client request
cat << EOF > client-csr.json
{
  "CN": "grpc.client",
  "hosts": [
    "grpc.client",
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C":  "CN",
      "L":  "XA",
      "ST": "XA",
      "O":  "grpc",
      "OU": "grpc"
    }
  ]
}
EOF

# client certificate
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile client client-csr.json | cfssljson -bare client
echo "sign CA client certificate successfully"

cat ca.pem >> client.pem

# clean
rm *.json *.csr
