import http from "utils/api/http";
import { SuccessResponse } from "utils/api/types";

export interface Gallery {
  title: string;
  id?: string;
  images?: string[];
  user_id?: number;
}

export type GalleryResponse = {
  gallery: Gallery;
};

export type GalleriesResponse = {
  galleries: Gallery[];
};

export async function createGallery(data: Gallery): Promise<GalleryResponse> {
  return http.post("/galleries", data);
}

export async function getGallery(id: string): Promise<GalleryResponse> {
  return http.get(`/galleries/${id}`);
}

export async function updateGallery(data: Gallery): Promise<GalleryResponse> {
  return http.post(`/galleries/${data.id}`, data);
}

export async function getAllGalleries(): Promise<GalleriesResponse> {
  return http.get("/galleries");
}

export async function deleteGallery(
  galleryId: string,
): Promise<SuccessResponse> {
  return http.post(`/galleries/${galleryId}/delete`);
}
