import { jsonPost } from "utils/api/fetch";

export type SignupFormData = {
  email: string;
  password: string;
};

export async function signUp(data: SignupFormData) {
  return jsonPost("/users", data);
}
