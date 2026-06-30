#include "../src/app.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

bool create_session(int user_id, char *token, size_t token_len) {
  char user_id_text[32];
  random_hex(token, 32);
  snprintf(user_id_text, sizeof(user_id_text), "%d", user_id);
  const char *params[2] = {token, user_id_text};
  PGresult *res = PQexecParams(
    db,
    "INSERT INTO sessions(token, user_id, expires_at) VALUES($1, $2::integer, now() + interval '7 days')",
    2,
    NULL,
    params,
    NULL,
    NULL,
    0
  );
  bool ok = PQresultStatus(res) == PGRES_COMMAND_OK;
  PQclear(res);
  if (!ok && token_len > 0) token[0] = '\0';
  return ok;
}

static bool extract_session_cookie(const char *cookie, char *token, size_t token_len) {
  const char *name = "sealion_session=";
  const char *p = strstr(cookie, name);
  if (!p) return false;
  p += strlen(name);
  size_t i = 0;
  while (p[i] && p[i] != ';' && i + 1 < token_len) {
    token[i] = p[i];
    i++;
  }
  token[i] = '\0';
  return i > 0;
}

bool current_user(const Request *req, User *user) {
  char token[129];
  if (!extract_session_cookie(req->cookie, token, sizeof(token))) {
    return false;
  }
  const char *params[1] = {token};
  PGresult *res = PQexecParams(
    db,
    "SELECT users.id, users.email "
    "FROM sessions JOIN users ON users.id = sessions.user_id "
    "WHERE sessions.token = $1 AND sessions.expires_at > now()",
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
  user->id = atoi(PQgetvalue(res, 0, 0));
  snprintf(user->email, sizeof(user->email), "%s", PQgetvalue(res, 0, 1));
  PQclear(res);
  return true;
}

void destroy_session(const char *cookie) {
  char token[129];
  if (!extract_session_cookie(cookie, token, sizeof(token))) return;
  const char *params[1] = {token};
  PGresult *res = PQexecParams(db, "DELETE FROM sessions WHERE token = $1", 1, NULL, params, NULL, NULL, 0);
  PQclear(res);
}
