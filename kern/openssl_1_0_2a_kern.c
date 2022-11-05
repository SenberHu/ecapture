#ifndef ECAPTURE_OPENSSL_1_0_2_A_KERN_H
#define ECAPTURE_OPENSSL_1_0_2_A_KERN_H

/* OPENSSL_VERSION_TEXT: OpenSSL 1.0.2u  20 Dec 2019, OPENSSL_VERSION_NUMBER: 268443999 */

// ssl_st->version
#define SSL_ST_VERSION 0x0

// ssl_st->session
#define SSL_ST_SESSION 0x130

// ssl_st->s3
#define SSL_ST_S3 0x80

// ssl_session_st->master_key
#define SSL_SESSION_ST_MASTER_KEY 0x14

// ssl3_state_st->client_random
#define SSL3_STATE_ST_CLIENT_RANDOM 0xc4

// ssl_session_st->cipher
#define SSL_SESSION_ST_CIPHER 0xe0

// ssl_session_st->cipher_id
#define SSL_SESSION_ST_CIPHER_ID 0xe8

// ssl_cipher_st->id
#define SSL_CIPHER_ST_ID 0x10

// openssl 1.0.2 does not support TLS 1.3, set 0 default
#define SSL_ST_HANDSHAKE_SECRET 0
#define SSL_ST_MASTER_SECRET 0
#define SSL_ST_SERVER_FINISHED_HASH 0
#define SSL_ST_HANDSHAKE_TRAFFIC_HASH 0
#define SSL_ST_EXPORTER_MASTER_SECRET 0

#include "openssl.h"
#include "openssl_masterkey.h"

#endif
