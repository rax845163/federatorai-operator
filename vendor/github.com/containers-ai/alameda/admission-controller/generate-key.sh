#!/bin/sh

NAMESPACE=${namespace:=alameda}
TMP_DIR="/tmp/admission-controller-certs"

mkdir -p ${TMP_DIR}

openssl genrsa -out ${TMP_DIR}/caKey.pem 2048
openssl req -x509 -new -nodes -key ${TMP_DIR}/caKey.pem -days 100000 -out ${TMP_DIR}/caCert.pem -subj "/CN=admission_ca"

cat > ${TMP_DIR}/server.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF

openssl genrsa -out ${TMP_DIR}/serverKey.pem 2048
openssl req -new -key ${TMP_DIR}/serverKey.pem -out ${TMP_DIR}/server.csr -subj "/CN=admission-controller.${NAMESPACE}.svc" -config ${TMP_DIR}/server.conf
openssl x509 -req -in ${TMP_DIR}/server.csr -CA ${TMP_DIR}/caCert.pem -CAkey ${TMP_DIR}/caKey.pem -CAcreateserial -out ${TMP_DIR}/serverCert.pem -days 100000 -extensions v3_req -extfile ${TMP_DIR}/server.conf

if [ -z "$(kubectl get ns | grep ${NAMESPACE})" ]
then 
kubectl create ns ${NAMESPACE} 
fi
kubectl create secret --namespace=${NAMESPACE} generic admission-controller-tls-certs --from-file=${TMP_DIR}/caKey.pem --from-file=${TMP_DIR}/caCert.pem --from-file=${TMP_DIR}/serverKey.pem --from-file=${TMP_DIR}/serverCert.pem
 
rm -rf ${TMP_DIR}
