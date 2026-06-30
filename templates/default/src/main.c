#define _GNU_SOURCE

#include <arpa/inet.h>
#include <ctype.h>
#include <errno.h>
#include <libpq-fe.h>
#include <netinet/in.h>
#include <openssl/rand.h>
#include <openssl/sha.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <time.h>
#include <unistd.h>

#define APP_NAME "__PROJECT_NAME__"
#define MAX_REQUEST 65536
#define MAX_BODY 32768
#define MAX_COOKIE 4096
#define MAX_PATH 1024
#define MAX_VIEW 32768

typedef struct {
  char method[8];
  char path[MAX_PATH];
  char cookie[MAX_COOKIE];
  char *body;
} Request;

typedef struct {
  int id;
  char email[256];
} User;

typedef struct {
  const char *name;
  const char *value;
} ViewVar;

static PGconn *db = NULL;

static void fatal(const char *message) {
  fprintf(stderr, "%s\n", message);
  exit(1);
}

static void send_all(int client, const char *data, size_t len) {
  size_t sent = 0;
  while (sent < len) {
    ssize_t n = send(client, data + sent, len - sent, 0);
    if (n <= 0) {
      return;
    }
    sent += (size_t)n;
  }
}

static void respond(int client, const char *status, const char *headers, const char *body) {
  char head[2048];
  size_t body_len = strlen(body);
  int n = snprintf(
    head,
    sizeof(head),
    "HTTP/1.1 %s\r\n"
    "Content-Type: text/html; charset=utf-8\r\n"
    "Content-Length: %zu\r\n"
    "Connection: close\r\n"
    "%s"
    "\r\n",
    status,
    body_len,
    headers ? headers : ""
  );
  send_all(client, head, (size_t)n);
  send_all(client, body, body_len);
}

static void redirect_to(int client, const char *location, const char *extra_headers) {
  char headers[1024];
  char body[1024];
  snprintf(
    headers,
    sizeof(headers),
    "Location: %s\r\n"
    "Cache-Control: no-store\r\n"
    "%s",
    location,
    extra_headers ? extra_headers : ""
  );
  snprintf(
    body,
    sizeof(body),
    "<!doctype html><title>Redirecting</title>"
    "<meta http-equiv=\"refresh\" content=\"0;url=%s\">"
    "<script>window.location.replace('%s');</script>"
    "<p><a href=\"%s\">Continue</a></p>",
    location,
    location,
    location
  );
  respond(client, "302 Found", headers, body);
}

static bool append_bytes(char *out, size_t out_len, size_t *used, const char *data, size_t data_len) {
  if (*used >= out_len) return false;
  size_t remaining = out_len - *used - 1;
  bool ok = data_len <= remaining;
  size_t copy_len = ok ? data_len : remaining;
  if (copy_len > 0) {
    memcpy(out + *used, data, copy_len);
    *used += copy_len;
  }
  out[*used] = '\0';
  return ok;
}

static bool append_text(char *out, size_t out_len, size_t *used, const char *text) {
  return append_bytes(out, out_len, used, text ? text : "", strlen(text ? text : ""));
}

static bool append_escaped(char *out, size_t out_len, size_t *used, const char *text) {
  const char *p = text ? text : "";
  while (*p) {
    switch (*p) {
      case '&':
        if (!append_text(out, out_len, used, "&amp;")) return false;
        break;
      case '<':
        if (!append_text(out, out_len, used, "&lt;")) return false;
        break;
      case '>':
        if (!append_text(out, out_len, used, "&gt;")) return false;
        break;
      case '"':
        if (!append_text(out, out_len, used, "&quot;")) return false;
        break;
      case '\'':
        if (!append_text(out, out_len, used, "&#39;")) return false;
        break;
      default:
        if (!append_bytes(out, out_len, used, p, 1)) return false;
        break;
    }
    p++;
  }
  return true;
}

static const char *view_var_value(const ViewVar *vars, size_t var_count, const char *name) {
  for (size_t i = 0; i < var_count; i++) {
    if (strcmp(vars[i].name, name) == 0) {
      return vars[i].value ? vars[i].value : "";
    }
  }
  return "";
}

static void copy_trimmed_key(char *out, size_t out_len, const char *start, size_t len) {
  while (len > 0 && isspace((unsigned char)*start)) {
    start++;
    len--;
  }
  while (len > 0 && isspace((unsigned char)start[len - 1])) {
    len--;
  }
  if (len >= out_len) len = out_len - 1;
  memcpy(out, start, len);
  out[len] = '\0';
}

static bool read_view_file(const char *path, char *out, size_t out_len) {
  FILE *file = fopen(path, "rb");
  if (!file) return false;
  size_t n = fread(out, 1, out_len - 1, file);
  out[n] = '\0';
  bool ok = feof(file);
  fclose(file);
  return ok;
}

static bool render_template_text(
  const char *template_text,
  const ViewVar *vars,
  size_t var_count,
  char *out,
  size_t out_len
) {
  const char *cursor = template_text;
  size_t used = 0;
  out[0] = '\0';

  while (*cursor) {
    const char *escaped = strstr(cursor, "{{");
    const char *raw = strstr(cursor, "{!!");
    bool use_raw = raw && (!escaped || raw < escaped);
    const char *start = use_raw ? raw : escaped;
    if (!start) {
      return append_text(out, out_len, &used, cursor);
    }

    if (!append_bytes(out, out_len, &used, cursor, (size_t)(start - cursor))) return false;

    const char *token_start = start + (use_raw ? 3 : 2);
    const char *end = use_raw ? strstr(token_start, "!!}") : strstr(token_start, "}}");
    if (!end) {
      return append_text(out, out_len, &used, start);
    }

    char key[128];
    copy_trimmed_key(key, sizeof(key), token_start, (size_t)(end - token_start));
    const char *value = view_var_value(vars, var_count, key);
    if (use_raw) {
      if (!append_text(out, out_len, &used, value)) return false;
    } else {
      if (!append_escaped(out, out_len, &used, value)) return false;
    }

    cursor = end + (use_raw ? 3 : 2);
  }

  return true;
}

static bool render_view_file(
  const char *path,
  const ViewVar *vars,
  size_t var_count,
  char *out,
  size_t out_len
) {
  char template_text[MAX_VIEW];
  if (!read_view_file(path, template_text, sizeof(template_text))) {
    return false;
  }
  return render_template_text(template_text, vars, var_count, out, out_len);
}

static bool render_page(
  const char *view_name,
  const char *title,
  const ViewVar *vars,
  size_t var_count,
  char *out,
  size_t out_len
) {
  char view_path[256];
  char content[MAX_VIEW];
  snprintf(view_path, sizeof(view_path), "view/%s.html", view_name);
  if (!render_view_file(view_path, vars, var_count, content, sizeof(content))) {
    return false;
  }

  ViewVar layout_vars[] = {
    {"title", title},
    {"app_name", APP_NAME},
    {"content", content},
  };
  return render_view_file("view/layout.html", layout_vars, 3, out, out_len);
}

static void respond_view(
  int client,
  const char *status,
  const char *view_name,
  const char *title,
  const ViewVar *vars,
  size_t var_count
) {
  char page[MAX_VIEW];
  if (!render_page(view_name, title, vars, var_count, page, sizeof(page))) {
    respond(
      client,
      "500 Internal Server Error",
      NULL,
      "<!doctype html><title>Template error</title><p>Could not render the requested view.</p>"
    );
    return;
  }
  respond(client, status, NULL, page);
}

static int hex_value(char c) {
  if (c >= '0' && c <= '9') return c - '0';
  if (c >= 'a' && c <= 'f') return c - 'a' + 10;
  if (c >= 'A' && c <= 'F') return c - 'A' + 10;
  return -1;
}

static void url_decode(char *out, size_t out_len, const char *in, size_t in_len) {
  size_t o = 0;
  for (size_t i = 0; i < in_len && o + 1 < out_len; i++) {
    if (in[i] == '+' ) {
      out[o++] = ' ';
    } else if (in[i] == '%' && i + 2 < in_len) {
      int hi = hex_value(in[i + 1]);
      int lo = hex_value(in[i + 2]);
      if (hi >= 0 && lo >= 0) {
        out[o++] = (char)((hi << 4) | lo);
        i += 2;
      }
    } else {
      out[o++] = in[i];
    }
  }
  out[o] = '\0';
}

static bool form_value(const char *body, const char *key, char *out, size_t out_len) {
  size_t key_len = strlen(key);
  const char *p = body ? body : "";
  while (*p) {
    const char *end = strchr(p, '&');
    size_t pair_len = end ? (size_t)(end - p) : strlen(p);
    if (pair_len > key_len && strncmp(p, key, key_len) == 0 && p[key_len] == '=') {
      url_decode(out, out_len, p + key_len + 1, pair_len - key_len - 1);
      return true;
    }
    if (!end) break;
    p = end + 1;
  }
  out[0] = '\0';
  return false;
}

static void random_hex(char *out, size_t byte_count) {
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

static void seed_admin(void) {
  PGresult *count = PQexec(db, "SELECT count(*) FROM users");
  if (PQresultStatus(count) != PGRES_TUPLES_OK) {
    PQclear(count);
    return;
  }
  int user_count = atoi(PQgetvalue(count, 0, 0));
  PQclear(count);
  if (user_count > 0) return;

  char salt[33];
  char hash[65];
  random_hex(salt, 16);
  password_hash("password", salt, hash);
  const char *params[3] = {"admin@sealion.local", hash, salt};
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
    fprintf(stderr, "could not seed admin: %s\n", PQerrorMessage(db));
  }
  PQclear(res);
}

static void connect_db(void) {
  const char *database_url = getenv("DATABASE_URL");
  if (!database_url) database_url = "postgres://sealion:sealion@localhost:5432/sealion";

  for (int attempt = 1; attempt <= 30; attempt++) {
    db = PQconnectdb(database_url);
    if (PQstatus(db) == CONNECTION_OK) {
      ensure_schema();
      seed_admin();
      return;
    }
    fprintf(stderr, "waiting for postgres (%d/30): %s", attempt, PQerrorMessage(db));
    PQfinish(db);
    db = NULL;
    sleep(1);
  }

  fatal("could not connect to postgres");
}

static bool create_user(const char *email, const char *password, char *error, size_t error_len) {
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

static bool verify_user(const char *email, const char *password, int *user_id) {
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

static bool lookup_user_id(const char *email, int *user_id) {
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

static bool create_session(int user_id, char *token, size_t token_len) {
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

static bool current_user(const Request *req, User *user) {
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

static void destroy_session(const char *cookie) {
  char token[129];
  if (!extract_session_cookie(cookie, token, sizeof(token))) return;
  const char *params[1] = {token};
  PGresult *res = PQexecParams(db, "DELETE FROM sessions WHERE token = $1", 1, NULL, params, NULL, NULL, 0);
  PQclear(res);
}

static void handle_home(int client, const Request *req) {
  User user;
  if (current_user(req, &user)) {
    redirect_to(client, "/dashboard", NULL);
    return;
  }
  ViewVar vars[] = {{"app_name", APP_NAME}};
  respond_view(client, "200 OK", "home", "Welcome", vars, 1);
}

static void handle_register_form(int client, const char *error) {
  ViewVar vars[] = {{"error", error && error[0] ? error : ""}};
  respond_view(client, "200 OK", "register", "Register", vars, 1);
}

static void handle_login_form(int client, const char *error) {
  ViewVar vars[] = {{"error", error && error[0] ? error : ""}};
  respond_view(client, "200 OK", "login", "Login", vars, 1);
}

static void handle_register(int client, const Request *req) {
  char email[256];
  char password[256];
  char error[512];
  form_value(req->body, "email", email, sizeof(email));
  form_value(req->body, "password", password, sizeof(password));
  if (!create_user(email, password, error, sizeof(error))) {
    char error_html[768];
    snprintf(error_html, sizeof(error_html), "<p class=\"error\">%s</p>", error);
    handle_register_form(client, error_html);
    return;
  }
  int user_id = 0;
  lookup_user_id(email, &user_id);
  char token[129];
  if (!create_session(user_id, token, sizeof(token))) {
    handle_login_form(client, "<p class=\"error\">Account created, but login failed. Try logging in.</p>");
    return;
  }
  char cookie[512];
  snprintf(cookie, sizeof(cookie), "Set-Cookie: sealion_session=%s; HttpOnly; SameSite=Lax; Path=/; Max-Age=604800\r\n", token);
  redirect_to(client, "/dashboard", cookie);
}

static void handle_login(int client, const Request *req) {
  char email[256];
  char password[256];
  int user_id = 0;
  form_value(req->body, "email", email, sizeof(email));
  form_value(req->body, "password", password, sizeof(password));
  if (!verify_user(email, password, &user_id)) {
    handle_login_form(client, "<p class=\"error\">Email or password is incorrect.</p>");
    return;
  }
  char token[129];
  if (!create_session(user_id, token, sizeof(token))) {
    handle_login_form(client, "<p class=\"error\">Could not create a session.</p>");
    return;
  }
  char cookie[512];
  snprintf(cookie, sizeof(cookie), "Set-Cookie: sealion_session=%s; HttpOnly; SameSite=Lax; Path=/; Max-Age=604800\r\n", token);
  redirect_to(client, "/dashboard", cookie);
}

static void handle_dashboard(int client, const Request *req) {
  User user;
  if (!current_user(req, &user)) {
    redirect_to(client, "/login", NULL);
    return;
  }
  ViewVar vars[] = {{"user_email", user.email}};
  respond_view(client, "200 OK", "dashboard", "Dashboard", vars, 1);
}

static void handle_logout(int client, const Request *req) {
  destroy_session(req->cookie);
  redirect_to(client, "/", "Set-Cookie: sealion_session=deleted; HttpOnly; SameSite=Lax; Path=/; Max-Age=0\r\n");
}

static int parse_content_length(const char *headers) {
  const char *p = strcasestr(headers, "Content-Length:");
  if (!p) return 0;
  p += strlen("Content-Length:");
  while (*p == ' ') p++;
  return atoi(p);
}

static void parse_cookie_header(const char *headers, char *cookie, size_t cookie_len) {
  const char *p = strcasestr(headers, "\r\nCookie:");
  if (!p) {
    cookie[0] = '\0';
    return;
  }
  p += strlen("\r\nCookie:");
  while (*p == ' ') p++;
  const char *end = strstr(p, "\r\n");
  size_t len = end ? (size_t)(end - p) : strlen(p);
  if (len >= cookie_len) len = cookie_len - 1;
  memcpy(cookie, p, len);
  cookie[len] = '\0';
}

static bool read_request(int client, Request *req, char *buffer, size_t buffer_len) {
  size_t total = 0;
  int content_length = 0;
  char *header_end = NULL;
  memset(req, 0, sizeof(*req));
  while (total + 1 < buffer_len) {
    ssize_t n = recv(client, buffer + total, buffer_len - total - 1, 0);
    if (n <= 0) break;
    total += (size_t)n;
    buffer[total] = '\0';
    header_end = strstr(buffer, "\r\n\r\n");
    if (header_end) {
      content_length = parse_content_length(buffer);
      size_t header_len = (size_t)(header_end + 4 - buffer);
      if (total >= header_len + (size_t)content_length) break;
    }
  }
  if (!header_end) return false;

  sscanf(buffer, "%7s %1023s", req->method, req->path);
  char *query = strchr(req->path, '?');
  if (query) *query = '\0';
  parse_cookie_header(buffer, req->cookie, sizeof(req->cookie));
  req->body = header_end + 4;
  return true;
}

static void handle_client(int client) {
  char buffer[MAX_REQUEST];
  Request req;
  if (!read_request(client, &req, buffer, sizeof(buffer))) {
    close(client);
    return;
  }

  printf("%s %s\n", req.method, req.path);
  fflush(stdout);

  if (strcmp(req.method, "GET") == 0 && strcmp(req.path, "/health") == 0) {
    respond(client, "200 OK", "Content-Type: text/plain; charset=utf-8\r\n", "ok\n");
  } else if (strcmp(req.method, "GET") == 0 && strcmp(req.path, "/") == 0) {
    handle_home(client, &req);
  } else if (strcmp(req.method, "GET") == 0 && strcmp(req.path, "/register") == 0) {
    handle_register_form(client, "");
  } else if (strcmp(req.method, "POST") == 0 && strcmp(req.path, "/register") == 0) {
    handle_register(client, &req);
  } else if (strcmp(req.method, "GET") == 0 && strcmp(req.path, "/login") == 0) {
    handle_login_form(client, "");
  } else if (strcmp(req.method, "POST") == 0 && strcmp(req.path, "/login") == 0) {
    handle_login(client, &req);
  } else if (strcmp(req.method, "GET") == 0 && strcmp(req.path, "/dashboard") == 0) {
    handle_dashboard(client, &req);
  } else if (strcmp(req.method, "POST") == 0 && strcmp(req.path, "/logout") == 0) {
    handle_logout(client, &req);
  } else {
    respond_view(client, "404 Not Found", "not_found", "Not found", NULL, 0);
  }

  close(client);
}

int main(void) {
  const char *port_text = getenv("APP_PORT");
  const char *public_url = getenv("PUBLIC_URL");
  int port = port_text ? atoi(port_text) : 8080;
  if (port <= 0) port = 8080;
  if (public_url && public_url[0] == '\0') public_url = NULL;

  connect_db();

  int server = socket(AF_INET, SOCK_STREAM, 0);
  if (server < 0) fatal("could not create socket");

  int yes = 1;
  setsockopt(server, SOL_SOCKET, SO_REUSEADDR, &yes, sizeof(yes));

  struct sockaddr_in addr;
  memset(&addr, 0, sizeof(addr));
  addr.sin_family = AF_INET;
  addr.sin_addr.s_addr = htonl(INADDR_ANY);
  addr.sin_port = htons((uint16_t)port);

  if (bind(server, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
    perror("bind");
    return 1;
  }
  if (listen(server, 64) < 0) {
    perror("listen");
    return 1;
  }

  if (public_url) {
    printf("%s listening inside container on :%d\n", APP_NAME, port);
    printf("open %s\n", public_url);
  } else {
    printf("%s listening on http://localhost:%d\n", APP_NAME, port);
  }
  fflush(stdout);

  for (;;) {
    int client = accept(server, NULL, NULL);
    if (client < 0) {
      if (errno == EINTR) continue;
      perror("accept");
      continue;
    }
    handle_client(client);
  }
}
