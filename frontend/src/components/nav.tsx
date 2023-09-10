import { NavLink } from "react-router-dom";
import useUserContext from "auth/user-provider";

export default function Nav() {
  const { handleSignOut, isSigningOut, user } = useUserContext();
  return (
    <header className="bg-gradient-to-r from-blue-800 to-indigo-800 text-white">
      <nav className="px-8 py-6 flex items-center space-x-12">
        <div className="text-4xl font-serif">Lenslocked</div>
        <div className="flex-grow">
          <NavLink
            className="text-lg font-semibold hover:text-blue-100 pr-8"
            to="/"
          >
            Home
          </NavLink>
          <NavLink
            className="text-lg font-semibold hover:text-blue-100 pr-8"
            to="/contact"
          >
            Contact
          </NavLink>
          <NavLink
            className="text-lg font-semibold hover:text-blue-100 pr-8"
            to="/faq"
          >
            FAQ
          </NavLink>
        </div>
        <div className="space-x-4">
          {user ? (
            <button onClick={handleSignOut} disabled={isSigningOut}>
              Sign out
            </button>
          ) : (
            <>
              <NavLink to="/signin">Sign in</NavLink>
              <NavLink
                to="/signup"
                className="px-4 py-2 bg-blue-700 hover:bg-blue-600 rounded"
              >
                Sign up
              </NavLink>
            </>
          )}
        </div>
      </nav>
    </header>
  );
}
