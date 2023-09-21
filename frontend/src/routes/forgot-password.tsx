import { useMutation } from "@tanstack/react-query";
import { Link, Navigate } from "react-router-dom";
import { Formik, Form, Field } from "formik";
import { ForgotPasswordFormData, requestForgotPassword } from "api/user";
import { useState } from "react";
import useUserContext from "auth/user-provider";

export default function ForgotPassword() {
  const { user } = useUserContext();
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [wasResetEmailSent, setWasResetEmailSent] = useState<boolean>(false);

  // Mutations
  const mutation = useMutation({
    mutationFn: requestForgotPassword,
    onSuccess: (res) => {
      if (res.success === true) {
        setWasResetEmailSent(true);
      }

      console.log(res);
      // setUser(res.user);
      setIsSubmitting(false);
    },
    onError: (e) => {
      console.log(e);
      setIsSubmitting(false);
    },
  });

  if (user) {
    return <Navigate to="/" replace />;
  }

  // computed
  const initialValues: ForgotPasswordFormData = {
    email: "",
  };

  return (
    <div className="py-12 flex justify-center">
      <div className="px-8 py-8 bg-white rounded shadow">
        {wasResetEmailSent ? (
          <>
            <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
              Check your email
            </h1>
            <p className="text-sm text-gray-600 pb-4">{`An email has been sent to the email address ${"hard coded"} with instructions to reset your password.`}</p>
          </>
        ) : (
          <>
            <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
              Forgot your password?
            </h1>
            <p className="text-sm text-gray-600 pb-4">
              No problem. Enter your email address and we'll send you a link to
              reset your password.
            </p>
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
                <div className="py-4">
                  <button
                    className="w-full py-4 px-2 bg-indigo-600 hover:bg-indigo-700
          text-white rounded font-bold text-lg"
                    disabled={isSubmitting}
                    type="submit"
                  >
                    Reset password
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
                    <Link to="/signin" className="underline">
                      Remember your password?
                    </Link>
                  </p>
                </div>
              </Form>
            </Formik>
          </>
        )}
      </div>
    </div>
  );
}
