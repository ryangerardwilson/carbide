#include "../src/app.h"

void handle_home(int client, const Request *req) {
  User user;
  if (current_user(req, &user)) {
    redirect_to(client, "/dashboard", NULL);
    return;
  }
  ViewVar vars[] = {{"app_name", APP_NAME}};
  respond_view(client, "200 OK", "home", "Welcome", vars, 1);
}

void handle_dashboard(int client, const Request *req) {
  User user;
  if (!current_user(req, &user)) {
    redirect_to(client, "/login", NULL);
    return;
  }
  ViewVar vars[] = {{"user_email", user.email}};
  respond_view(client, "200 OK", "dashboard", "Dashboard", vars, 1);
}

void handle_not_found(int client) {
  respond_view(client, "404 Not Found", "not_found", "Not found", NULL, 0);
}
