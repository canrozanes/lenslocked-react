import { User } from "models/user";
import http from "utils/api/http";
import { SuccessResponse } from "utils/api/types";

export type SignupFormData = {
  email: string;
  password: string;
};

type UserResponse = {
  user: User;
};

export async function signUp(data: SignupFormData): Promise<UserResponse> {
  return http.post("/users", data);
}

export async function signIn(data: SignupFormData): Promise<UserResponse> {
  return http.post("/signin", data);
}

export async function getMe(): Promise<UserResponse> {
  return http.get("/users/me");
}

export async function signOut(): Promise<SuccessResponse> {
  return http.post("/signout");
}
