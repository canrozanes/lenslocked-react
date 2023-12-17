import { useQuery } from "@tanstack/react-query";
import { getAllGalleries } from "api/gallery";
import GalleriesTable from "components/gallery/galleries-table";
import { NavLink } from "react-router-dom";

export default function GalleriesIndex() {
  const getAllGalleriesQuery = useQuery({
    queryFn: getAllGalleries,
    queryKey: ["galleries"],
  });

  if (!getAllGalleriesQuery.data) {
    return "loading";
  }

  return (
    <div className="p-8 w-full">
      <h1 className="pt-4 pb-8 text-3xl font-bold text-gray-800">
        My Galleries
      </h1>
      <GalleriesTable galleries={getAllGalleriesQuery.data.galleries} />
      <div className="py-4">
        <NavLink
          to="/galleries/new"
          className="py-2 px-8 bg-indigo-600 hover:bg-indigo-700 text-lg text-white font-bold rounded"
        >
          New Gallery
        </NavLink>
      </div>
    </div>
  );
}
