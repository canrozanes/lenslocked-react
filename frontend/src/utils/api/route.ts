const enum Route {
  Development = "http://localhost:3000/api",
  Production = "TODO",
}

function buildUrl() {
  if (import.meta.env.MODE === "development") {
    return Route.Development;
  }

  return Route.Production;
}

export const API_ROUTE = buildUrl();
