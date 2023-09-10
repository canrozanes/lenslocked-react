import http from "utils/api/http";
import { SuccessResponse } from "utils/api/types";

export async function getCsrf(): Promise<SuccessResponse> {
  return http.get("/csrf");
}
