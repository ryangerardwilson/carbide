#ifndef SEALION_APP_H
#define SEALION_APP_H

#include <libpq-fe.h>
#include <stdbool.h>
#include <stddef.h>

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

extern PGconn *db;

void fatal(const char *message);
void respond(int client, const char *status, const char *headers, const char *body);
void redirect_to(int client, const char *location, const char *extra_headers);
void respond_view(
  int client,
  const char *status,
  const char *view_name,
  const char *title,
  const ViewVar *vars,
  size_t var_count
);
bool form_value(const char *body, const char *key, char *out, size_t out_len);

void connect_db(void);
void random_hex(char *out, size_t byte_count);
bool create_user(const char *email, const char *password, char *error, size_t error_len);
bool verify_user(const char *email, const char *password, int *user_id);
bool lookup_user_id(const char *email, int *user_id);
bool create_session(int user_id, char *token, size_t token_len);
bool current_user(const Request *req, User *user);
void destroy_session(const char *cookie);

void handle_home(int client, const Request *req);
void handle_register_form(int client, const char *error);
void handle_login_form(int client, const char *error);
void handle_register(int client, const Request *req);
void handle_login(int client, const Request *req);
void handle_dashboard(int client, const Request *req);
void handle_logout(int client, const Request *req);
void handle_not_found(int client);

#endif
