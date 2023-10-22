import { useMutation } from "@tanstack/react-query";
import { Link, Navigate } from "react-router-dom";
import { Formik, Form, Field } from "formik";
import { SignupFormData, signIn } from "api/user";
import { useState } from "react";
import useUserContext from "auth/user-provider";
import useAlert from "alerts/alert-context";
import { AxiosError } from "axios";

export default function SignIn() {
  const { user, setUser } = useUserContext();
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const { setAlert } = useAlert();

  // Mutations
  const mutation = useMutation({
    mutationFn: signIn,
    onSuccess: (res) => {
      setUser(res.user);
      setIsSubmitting(false);
    },
    onError: (e: AxiosError<{ data: string }>) => {
      if (e.response?.status === 401) {
        setAlert("invalid email/password");
      } else {
        setAlert("Something went wrong please try again");
      }
      setIsSubmitting(false);
    },
  });

  if (user) {
    return <Navigate to="/" replace />;
  }

  // computed
  const initialValues: SignupFormData = {
    email: "",
    password: "",
  };

  return (
    <div className="py-12 flex justify-center">
      <div className="px-8 py-8 bg-white rounded shadow">
        <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
          Welcome back!
        </h1>
        <Formik
          initialValues={initialValues}
          onSubmit={(values) => {
            mutation.reset();
            setIsSubmitting(true);
            mutation.mutate(values);
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
                Sign in
              </button>
            </div>
            <div className="py-2 w-full flex justify-between">
              <p className="text-xs text-gray-500">
                Need an account?
                <Link to="/signup" className="underline">
                  Sign up
                </Link>
              </p>
              <p className="text-xs text-gray-500">
                <Link to="/forgot-pw" className="underline">
                  Forgot your password?
                </Link>
              </p>
            </div>
          </Form>
        </Formik>
      </div>
    </div>
  );
}
