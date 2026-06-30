#define _GNU_SOURCE

#include "../src/app.h"

#include <openssl/rand.h>
#include <openssl/sha.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

void random_hex(char *out, size_t byte_count) {
  unsigned char bytes[64];
  if (byte_count > sizeof(bytes)) byte_count = sizeof(bytes);
  if (RAND_bytes(bytes, (int)byte_count) != 1) {
    srand((unsigned int)time(NULL));
    for (size_t i = 0; i < byte_count; i++) {
      bytes[i] = (unsigned char)(rand() & 0xff);
    }
  }
  for (size_t i = 0; i < byte_count; i++) {
    sprintf(out + (i * 2), "%02x", bytes[i]);
  }
  out[byte_count * 2] = '\0';
}

static void password_hash(const char *password, const char *salt, char *out) {
  char input[1024];
  unsigned char digest[SHA256_DIGEST_LENGTH];
  snprintf(input, sizeof(input), "%s:%s", salt, password);
  SHA256((const unsigned char *)input, strlen(input), digest);
  for (size_t i = 0; i < SHA256_DIGEST_LENGTH; i++) {
    sprintf(out + (i * 2), "%02x", digest[i]);
  }
  out[SHA256_DIGEST_LENGTH * 2] = '\0';
}

static void db_exec_or_die(const char *sql) {
  PGresult *res = PQexec(db, sql);
  ExecStatusType status = PQresultStatus(res);
  if (status != PGRES_COMMAND_OK && status != PGRES_TUPLES_OK) {
    fprintf(stderr, "database error: %s\nsql: %s\n", PQerrorMessage(db), sql);
    PQclear(res);
    exit(1);
  }
  PQclear(res);
}

static void ensure_schema(void) {
  db_exec_or_die(
    "CREATE TABLE IF NOT EXISTS users ("
    "id SERIAL PRIMARY KEY,"
    "email TEXT NOT NULL UNIQUE,"
    "password_hash TEXT NOT NULL,"
    "password_salt TEXT NOT NULL,"
    "created_at TIMESTAMPTZ NOT NULL DEFAULT now()"
    ")"
  );
  db_exec_or_die(
    "CREATE TABLE IF NOT EXISTS sessions ("
    "token TEXT PRIMARY KEY,"
    "user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,"
    "created_at TIMESTAMPTZ NOT NULL DEFAULT now(),"
    "expires_at TIMESTAMPTZ NOT NULL"
    ")"
  );
  db_exec_or_die("CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions(user_id)");
  db_exec_or_die("CREATE INDEX IF NOT EXISTS sessions_expires_at_idx ON sessions(expires_at)");
  db_exec_or_die("DELETE FROM sessions WHERE expires_at < now()");
}

void connect_db(void) {
  const char *database_url = getenv("DATABASE_URL");
  if (!database_url) database_url = "postgres://sealion:sealion@localhost:5432/sealion";

  for (int attempt = 1; attempt <= 30; attempt++) {
    db = PQconnectdb(database_url);
    if (PQstatus(db) == CONNECTION_OK) {
      ensure_schema();
      return;
    }
    fprintf(stderr, "waiting for postgres (%d/30): %s", attempt, PQerrorMessage(db));
    PQfinish(db);
    db = NULL;
    sleep(1);
  }

  fatal("could not connect to postgres");
}

bool create_user(const char *email, const char *password, char *error, size_t error_len) {
  if (strlen(email) < 3 || strchr(email, '@') == NULL) {
    snprintf(error, error_len, "Enter a valid email address.");
    return false;
  }
  if (strlen(password) < 6) {
    snprintf(error, error_len, "Password must be at least 6 characters.");
    return false;
  }

  char salt[33];
  char hash[65];
  random_hex(salt, 16);
  password_hash(password, salt, hash);
  const char *params[3] = {email, hash, salt};
  PGresult *res = PQexecParams(
    db,
    "INSERT INTO users(email, password_hash, password_salt) VALUES($1, $2, $3)",
    3,
    NULL,
    params,
    NULL,
    NULL,
    0
  );
  if (PQresultStatus(res) != PGRES_COMMAND_OK) {
    snprintf(error, error_len, "That email is already registered.");
    PQclear(res);
    return false;
  }
  PQclear(res);
  error[0] = '\0';
  return true;
}

bool verify_user(const char *email, const char *password, int *user_id) {
  const char *params[1] = {email};
  PGresult *res = PQexecParams(
    db,
    "SELECT id, password_hash, password_salt FROM users WHERE email = $1",
    1,
    NULL,
    params,
    NULL,
    NULL,
    0
  );
  if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) != 1) {
    PQclear(res);
    return false;
  }
  char hash[65];
  password_hash(password, PQgetvalue(res, 0, 2), hash);
  bool ok = strcmp(hash, PQgetvalue(res, 0, 1)) == 0;
  if (ok) *user_id = atoi(PQgetvalue(res, 0, 0));
  PQclear(res);
  return ok;
}

bool lookup_user_id(const char *email, int *user_id) {
  const char *params[1] = {email};
  PGresult *res = PQexecParams(db, "SELECT id FROM users WHERE email = $1", 1, NULL, params, NULL, NULL, 0);
  if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) != 1) {
    PQclear(res);
    return false;
  }
  *user_id = atoi(PQgetvalue(res, 0, 0));
  PQclear(res);
  return true;
}
