export default function SignUp() {
  return (
    <div className="py-12 flex justify-center">
      <div className="px-8 py-8 bg-white rounded shadow">
        <h1 className="pt-4 pb-8 text-center text-3xl font-bold text-gray-900">
          Start sharing your photos today!
        </h1>
        <form action="/signup" method="post">
          <div className="py-2">
            <label
              htmlFor="email"
              className="text-sm font-semibold text-gray-800"
            >
              Email Address
            </label>
            <input
              name="email"
              id="email"
              type="email"
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
            <input
              name="password"
              id="password"
              type="password"
              placeholder="Password"
              required
              className="w-full px-3 py-2 border border-gray-300 placeholder-gray-500
            text-gray-800 rounded"
            />
          </div>
          <div className="py-4">
            <button
              className="w-full py-4 px-2 bg-indigo-600 hover:bg-indigo-700
          text-white rounded font-bold text-lg"
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
        </form>
      </div>
    </div>
  );
}
