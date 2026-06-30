#include "../src/app.h"

#include <stdio.h>

void handle_register_form(int client, const char *error) {
  ViewVar vars[] = {
    {"auth_title", "Create account"},
    {"auth_action", "/register"},
    {"email_value", ""},
    {"password_autocomplete", "new-password"},
    {"submit_label", "Create account"},
    {"error", error && error[0] ? error : ""},
    {"auth_footer", "<p class=\"muted\">Already have an account? <a href=\"/login\">Log in</a>.</p>"},
  };
  respond_view(client, "200 OK", "register", "Register", vars, 7);
}

void handle_login_form(int client, const char *error) {
  ViewVar vars[] = {
    {"auth_title", "Log in"},
    {"auth_action", "/login"},
    {"email_value", "admin@sealion.local"},
    {"password_autocomplete", "current-password"},
    {"submit_label", "Log in"},
    {"error", error && error[0] ? error : ""},
    {"auth_footer", "<p class=\"muted\">Demo password: <code>password</code></p>"},
  };
  respond_view(client, "200 OK", "login", "Login", vars, 7);
}

void handle_register(int client, const Request *req) {
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

void handle_login(int client, const Request *req) {
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

void handle_logout(int client, const Request *req) {
  destroy_session(req->cookie);
  redirect_to(client, "/", "Set-Cookie: sealion_session=deleted; HttpOnly; SameSite=Lax; Path=/; Max-Age=0\r\n");
}
