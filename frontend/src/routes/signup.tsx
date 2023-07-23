import { useMutation } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { Formik, Form, Field } from "formik";
import { SignupFormData, signUp } from "api/user";
import { useState } from "react";

export default function SignUp() {
  const [searchParams] = useSearchParams();
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  // Mutations
  const { mutate } = useMutation({
    mutationFn: signUp,
    onSuccess: () => {
      setIsSubmitting(false);
    },
    onError: () => {},
  });

  // computed
  const initialValues: SignupFormData = {
    email: searchParams.get("email") ?? "",
    password: "",
  };

  return (
    <div className="py-12 flex justify-center">
      <div className="px-8 py-8 bg-white rounded shadow">
        <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
          Start sharing your photos today!
        </h1>
        <Formik
          initialValues={initialValues}
          onSubmit={(values) => {
            setIsSubmitting(true);
            mutate(values);
          }}
        >
          <Form>
            <div className="py-2">
              <label
                htmlFor="email"
                className="text-sm font-semibold text-gray-800"
              >
                Email Address
              </label>
              <Field
                type="email"
                name="email"
                id="email"
                placeholder="Email address"
                required
                autoComplete="email"
                className="w-full px-3 py-2 border border-gray-300 placeholder-gray-500
            text-gray-800 rounded"
              />
            </div>
            <div className="py-2">
              <label
                htmlFor="password"
                className="text-sm font-semibold text-gray-800"
              >
                Password
              </label>
              <Field
                type="password"
                name="password"
                id="password"
                placeholder="Password"
                required
                autoComplete="current-password"
                className="w-full px-3 py-2 border border-gray-300 placeholder-gray-500
            text-gray-800 rounded"
              />
            </div>
            <div className="py-4">
              <button
                className="w-full py-4 px-2 bg-indigo-600 hover:bg-indigo-700
          text-white rounded font-bold text-lg"
                disabled={isSubmitting}
                type="submit"
              >
                Sign up
              </button>
            </div>
            <div className="py-2 w-full flex justify-between">
              <p className="text-xs text-gray-500">
                Already have an account?
                <a href="/signin" className="underline">
                  Sign in
                </a>
              </p>
              <p className="text-xs text-gray-500">
                <a href="/reset-pw" className="underline">
                  Forgot your password?
                </a>
              </p>
            </div>
          </Form>
        </Formik>
      </div>
    </div>
  );
}
