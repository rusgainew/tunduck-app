import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  index("routes/_index.tsx"),
  route("login", "routes/login.tsx"),
  layout("routes/dashboard.tsx", [
    index("routes/dashboard/_index.tsx"),
    route("users", "routes/dashboard/users.tsx"),
    route("organizations", "routes/dashboard/organizations.tsx"),
    route("documents", "routes/dashboard/documents.tsx"),
    route("roles", "routes/dashboard/roles.tsx"),
  ]),
] satisfies RouteConfig;
