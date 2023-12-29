import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { Gallery, GalleryResponse, createGallery } from "api/gallery";
import useAlert from "alerts/alert-context";
import { AxiosError } from "axios";
import GalleryForm from "components/gallery/gallery-form";
import { useNavigate } from "react-router-dom";

export default function New() {
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const { setAlert } = useAlert();
  const navigate = useNavigate();

  // Mutations
  const { mutate } = useMutation({
    mutationFn: createGallery,
    onSuccess: (_: GalleryResponse) => {
      setIsSubmitting(false);
      navigate("/galleries");
    },
    onError: (e: AxiosError) => {
      console.error(e);
      setAlert("Something went wrong. Please try again");

      setIsSubmitting(false);
    },
  });

  // computed
  const initialValues: Gallery = {
    title: "",
  };

  const onSubmit = (values: Gallery) => {
    setAlert("");
    setIsSubmitting(true);
    mutate(values);
  };

  return (
    <div className="p-8 w-full">
      <h1 className="pt-4 pb-8 text-3xl font-bold text-gray-800">
        Create a new Gallery
      </h1>
      <GalleryForm
        isSubmitting={isSubmitting}
        initialValues={initialValues}
        onSubmit={onSubmit}
      />
    </div>
  );
}
