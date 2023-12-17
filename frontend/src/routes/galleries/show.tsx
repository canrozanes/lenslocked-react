import { useQuery } from "@tanstack/react-query";
import { GalleryResponse, getGallery } from "api/gallery";
import { AxiosError } from "axios";

import { useParams } from "react-router-dom";

export default function ShowGallery() {
  const params = useParams();
  const id = params.id ?? "";

  const getAllGalleriesQuery = useQuery<GalleryResponse, AxiosError>({
    queryFn: () => getGallery(id),
    queryKey: ["gallery"],
  });

  if (getAllGalleriesQuery.isFetching) {
    return "loading";
  }

  if (getAllGalleriesQuery.error || !getAllGalleriesQuery.data?.gallery) {
    return "error";
  }

  const gallery = getAllGalleriesQuery.data.gallery;

  return (
    <div className="px-8 py-12 w-full">
      <h1 className="pt-4 pb-8 text-3xl font-bold text-gray-900">
        {gallery?.title}
      </h1>
      <div className="columns-4 gap-4 space-y-4">
        {gallery?.images?.map((image) => (
          <div>
            <a href={image}>
              <img className="w-full" src={image} />
            </a>
          </div>
        ))}
      </div>
    </div>
  );
}
