import http from "utils/api/fetch";

export type SignupFormData = {
  email: string;
  password: string;
};

export async function signUp(data: SignupFormData) {
  return http.post("/users", data);
}

export async function signIn(data: SignupFormData) {
  return http.post("/signin", data);
}

export async function getMe() {
  return http.get("/users/me");
}
