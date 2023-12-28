import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import useAlert from "alerts/alert-context";
import {
  Gallery,
  GalleryResponse,
  deleteGallery,
  getGallery,
  updateGallery,
} from "api/gallery";
import { AxiosError } from "axios";
import GalleryForm from "components/gallery/gallery-form";
import { useNavigate, useParams } from "react-router-dom";
import useUserContext from "auth/user-provider";

export default function GalleriesEdit() {
  const params = useParams();
  const id = params.id ?? "";

  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const { setAlert } = useAlert();
  const { user } = useUserContext();
  const queryClient = useQueryClient();

  const navigate = useNavigate();

  const getGalleryQuery = useQuery<GalleryResponse, AxiosError>({
    queryFn: () => getGallery(id),
    queryKey: ["gallery"],
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

  const deleteMutation = useMutation({
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

  const onSubmit = (values: Gallery) => {
    setAlert("");
    setIsSubmitting(true);
    updateMutation.mutate(values);
  };

  if (getGalleryQuery.isFetching) {
    return <p>loading</p>;
  }

  if (user?.id && user.id !== getGalleryQuery.data?.gallery?.user_id) {
    return <p>You don't have edit access to this gallery</p>;
  }

  return (
    <div className="p-8 w-full">
      <h1 className="pt-4 pb-8 text-3xl font-bold text-gray-800">
        Edit Gallery
      </h1>
      {!getGalleryQuery.error && getGalleryQuery.data ? (
        <>
          <GalleryForm
            isSubmitting={isSubmitting}
            initialValues={getGalleryQuery.data.gallery}
            onSubmit={onSubmit}
          />
          <div className="py-4">
            <h2>Dangerus actions</h2>
            <button
              type="submit"
              className="py-2 px-8 bg-red-600 hover:bg-red-700 text-white rounded font-bold text-lg"
              onClick={() => {
                debugger;
                setAlert("");
                setIsSubmitting(true);
                deleteMutation.mutate(id);
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
