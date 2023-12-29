import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import useAlert from "alerts/alert-context";
import {
  Gallery,
  GalleryResponse,
  deleteGallery,
  getGallery,
  updateGallery,
  deleteImage,
  uploadImages,
  UploadImageVariables,
  DeleteImageVariables,
} from "api/gallery";
import { AxiosError } from "axios";
import GalleryForm from "components/gallery/gallery-form";
import { useNavigate, useParams } from "react-router-dom";
import useUserContext from "auth/user-provider";

export default function GalleriesEdit() {
  const params = useParams();
  const id = params.id ?? "";

  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [images, setImages] = useState<FileList | null>(null);
  const { setAlert } = useAlert();
  const { user } = useUserContext();
  const queryClient = useQueryClient();

  const navigate = useNavigate();

  const getGalleryQuery = useQuery<GalleryResponse, AxiosError>({
    queryFn: () => getGallery(id),
    queryKey: ["gallery"],
    refetchOnWindowFocus: false,
  });

  // Mutations
  const updateMutation = useMutation({
    mutationFn: updateGallery,
    onSuccess: () => {
      setIsSubmitting(false);
      setAlert("Gallery successfully deleted");
    },
    onError: (e: AxiosError) => {
      console.error(e);
      setAlert("Something went wrong. Please try again");

      setIsSubmitting(false);
    },
  });

  const deleteGalleryMutation = useMutation({
    mutationFn: deleteGallery,
    onSuccess: () => {
      setIsSubmitting(false);
      navigate("/galleries");
      setAlert("Gallery successfully deleted");
      queryClient.invalidateQueries({ queryKey: ["galleries"] });
    },
    onError: (e: AxiosError) => {
      setAlert("Something went wrong. Please try again");
      console.error(e);
      setIsSubmitting(false);
    },
  });

  const deleteImageMutation = useMutation({
    mutationFn: (variables: DeleteImageVariables) => deleteImage(variables),
    onSuccess: () => {
      setIsSubmitting(false);
      setAlert("Image successfully deleted");
      queryClient.invalidateQueries({ queryKey: ["gallery"] });
    },
    onError: (e: AxiosError) => {
      setAlert("Something went wrong. Please try again");
      console.error(e);
      setIsSubmitting(false);
    },
  });

  const uploadImageMutation = useMutation({
    mutationFn: (variables: UploadImageVariables) => uploadImages(variables),
    onSuccess: () => {
      setIsSubmitting(false);
      setAlert("Image successfully uploaded");
      queryClient.invalidateQueries({ queryKey: ["gallery"] });
    },
    onError: (e: AxiosError) => {
      setAlert("Something went wrong. Please try again");
      console.error(e);
      setIsSubmitting(false);
    },
  });

  const onSubmit = (values: Gallery) => {
    setAlert("");
    setIsSubmitting(true);
    updateMutation.mutate(values);
  };

  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    debugger;
    const files = e.target.files;
    setImages(files);
  };

  if (getGalleryQuery.isFetching) {
    return <p>loading</p>;
  }

  if (user?.id && user.id !== getGalleryQuery.data?.gallery?.user_id) {
    return <p>You don't have edit access to this gallery</p>;
  }

  const gallery = getGalleryQuery.data?.gallery;

  if (!gallery) {
    return <p>Gallery not found</p>;
  }

  return (
    <div className="p-8 w-full">
      <h1 className="pt-4 pb-8 text-3xl font-bold text-gray-800">
        Edit Gallery
      </h1>
      {!getGalleryQuery.error && gallery ? (
        <>
          <GalleryForm
            isSubmitting={isSubmitting}
            initialValues={gallery}
            onSubmit={onSubmit}
          />
          <div className="py-4">
            <div className="py-2">
              <label
                htmlFor="images"
                className="block mb-2 text-sm font-semibold text-gray-800"
              >
                Add Images
                <p className="py-2 text-xs text-gray-600 font-normal">
                  Please only upload jpg, png, and gif files.
                </p>
              </label>
              <input
                type="file"
                multiple
                accept="image/png, image/jpeg, image/gif"
                id="images"
                name="images"
                onChange={onFileChange}
              />
            </div>
            <button
              className="py-2 px-8 bg-indigo-600 hover:bg-indigo-700 text-white text-lg font-bold rounded"
              onClick={() => {
                debugger;
                if (!images) {
                  return;
                }
                setAlert("");
                setIsSubmitting(true);
                uploadImageMutation.mutate({ galleryId: id, files: images });
              }}
            >
              Upload
            </button>
          </div>
          <div className="py-4">
            <h2 className="pb-2 text-sm font-semibold text-gray-800">
              Current Images
            </h2>
            <div className="py-2 grid grid-cols-8 gap-2">
              {gallery.images?.map((image) => (
                <div
                  className="h-min w-full relative"
                  key={image.filename + image.gallery_id}
                >
                  <div className="absolute top-2 right-2">
                    <button
                      onClick={() => {
                        deleteImageMutation.mutate({
                          filename: image.filename,
                          galleryId: image.gallery_id,
                        });
                      }}
                      className="p-1 text-xs text-red-800 bg-red-100 border border-red-400 rounded"
                    >
                      Delete
                    </button>
                  </div>
                  <img
                    className="w-full"
                    src={`/api/galleries/${image.gallery_id}/images/${image.filename_escaped}`}
                  />
                </div>
              ))}
            </div>
          </div>
          <div className="py-4">
            <h2>Dangerus actions</h2>
            <button
              className="py-2 px-8 bg-red-600 hover:bg-red-700 text-white rounded font-bold text-lg"
              onClick={() => {
                setAlert("");
                setIsSubmitting(true);
                deleteGalleryMutation.mutate(id);
              }}
            >
              Delete
            </button>
          </div>
        </>
      ) : (
        <>loading</>
      )}
    </div>
  );
}
