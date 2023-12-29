import http from "utils/api/http";
import { SuccessResponse } from "utils/api/types";

export interface Image {
  id?: string;
  filename: string;
  filename_escaped: string;
  gallery_id: string;
}

export interface Gallery {
  title: string;
  id?: string;
  images?: Image[];
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

export type DeleteImageVariables = {
  galleryId: string;
  filename: string;
};

export async function deleteImage({
  galleryId,
  filename,
}: DeleteImageVariables): Promise<SuccessResponse> {
  return http.post(`/galleries/${galleryId}/images/${filename}/delete`);
}

export type UploadImageVariables = {
  galleryId: string;
  files: FileList;
};

export async function uploadImages({
  galleryId,
  files,
}: UploadImageVariables): Promise<SuccessResponse> {
  const formData = new FormData();
  for (let file of files) {
    formData.append("files", file);
  }
  const config = {
    headers: {
      "content-type": "multipart/form-data",
    },
  };

  return http.post(`/galleries/${galleryId}/images`, formData, config);
}
