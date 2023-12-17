import { Gallery } from "api/gallery";
import { Field, Form, Formik } from "formik";
import * as Yup from "yup";

type GalleryFormProps = {
  isSubmitting: boolean;
  initialValues: Gallery;
  onSubmit: (values: Gallery) => void;
  isNew?: boolean;
};

const GallerySchema = Yup.object().shape({
  title: Yup.string()
    .min(2, "Too Short!")
    .max(50, "Too Long!")
    .required("Required"),
});

export default function GalleryForm({
  isSubmitting,
  initialValues,
  onSubmit,
}: GalleryFormProps) {
  return (
    <Formik
      initialValues={initialValues}
      validationSchema={GallerySchema}
      onSubmit={onSubmit}
    >
      {({ errors, touched }) => (
        <Form>
          <div className="py-2">
            <label
              htmlFor="title"
              className="text-sm font-semibold text-gray-800"
            >
              Title
            </label>
            <Field
              name="title"
              id="title"
              type="text"
              placeholder="Gallery Title"
              required
              className="w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-800 rounded"
              // value="{{.Title}}"
              autoFocus
            />
            {errors.title && touched.title ? (
              <div className="text-red-700">{errors.title}</div>
            ) : null}
          </div>

          <div className="py-4">
            <button
              type="submit"
              className="py-2 px-8 bg-indigo-600 hover:bg-indigo-700 text-white rounded font-bold text-lg"
              disabled={isSubmitting}
            >
              Save
            </button>
          </div>
        </Form>
      )}
    </Formik>
  );
}
