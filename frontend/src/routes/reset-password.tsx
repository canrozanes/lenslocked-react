import { useMutation, useQuery } from "@tanstack/react-query";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router-dom";
import { Formik, Form, Field } from "formik";
import { ResetPasswordFormData, getMe, resetPassword } from "api/user";
import { useState } from "react";
import useUserContext from "auth/user-provider";

export default function ResetPassword() {
  const { user, setUser } = useUserContext();
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [isPasswordResetSuccess, setIsPasswordResetSuccess] =
    useState<boolean>(false);
  let [searchParams] = useSearchParams();
  const navigate = useNavigate();

  // Mutations
  const mutation = useMutation({
    mutationFn: resetPassword,
    onSuccess: () => {
      debugger;
      setIsPasswordResetSuccess(true);
      setIsSubmitting(false);
    },
    onError: (e) => {
      console.log(e);
      setIsSubmitting(false);
    },
  });

  useQuery({
    queryFn: getMe,
    enabled: isPasswordResetSuccess,
    onSuccess: (res) => {
      debugger;
      setUser(res.user);
    },
    onError: () => {
      navigate("/signin");
    },
  });

  if (user) {
    return <Navigate to="/" replace />;
  }

  // computed
  const initialValues: ResetPasswordFormData = {
    password: "",
    token: searchParams.get("token") ?? "",
  };

  return (
    <div className="py-12 flex justify-center">
      <div className="px-8 py-8 bg-white rounded shadow">
        <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
          Reset your password
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
                htmlFor="password"
                className="text-sm font-semibold text-gray-800"
              >
                New password
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
            <div className="hidden">
              <Field type="token" name="token" id="token" />
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
