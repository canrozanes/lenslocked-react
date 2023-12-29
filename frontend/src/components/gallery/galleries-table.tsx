import { useMutation, useQueryClient } from "@tanstack/react-query";
import useAlert from "alerts/alert-context";
import { Gallery, deleteGallery } from "api/gallery";
import { AxiosError } from "axios";
import { useState } from "react";
import { NavLink } from "react-router-dom";

type GalleriesTableProps = {
  galleries: Gallery[];
};

export default function GalleriesTable({ galleries }: GalleriesTableProps) {
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const { setAlert } = useAlert();

  const queryClient = useQueryClient();

  const { mutate } = useMutation({
    mutationFn: deleteGallery,
    onSuccess: () => {
      setIsSubmitting(false);
      queryClient.invalidateQueries({ queryKey: ["galleries"] });
    },
    onError: (e: AxiosError) => {
      setAlert("Something went wrong. Please try again");
      console.error(e);
      setIsSubmitting(false);
    },
  });
  return (
    <table className="w-full table-fixed">
      <thead>
        <tr>
          <th className="p-2 text-left w-24">ID</th>
          <th className="p-2 text-left">Title</th>
          <th className="p-2 text-left w-96">Actions</th>
        </tr>
      </thead>
      <tbody>
        {galleries.map(({ id, title }) => (
          <tr key={id} className="border">
            <td className="p-2 border">{id}</td>
            <td className="p-2 border">{title}</td>
            <td className="p-2 border flex space-x-2">
              <NavLink
                className="py-1 px-2 bg-blue-100 hover:bg-blue-200 rounded border border-blue-600 text-xs text-blue-600"
                to={`/galleries/${id}`}
              >
                View
              </NavLink>
              <NavLink
                className="py-1 px-2 bg-blue-100 hover:bg-blue-200 rounded border border-blue-600 text-xs text-blue-600"
                to={`/galleries/${id}/edit`}
              >
                Edit
              </NavLink>
              <button
                className="py-1 px-2 bg-red-100 hover:bg-red-200 rounded border border-red-600 text-xs text-red-600"
                onClick={() => {
                  setAlert("");
                  setIsSubmitting(true);
                  mutate(id!);
                }}
                disabled={isSubmitting}
              >
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
