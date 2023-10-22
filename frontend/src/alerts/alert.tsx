import useAlert from "alerts/alert-context";

export default function Alert() {
  const { text, setAlert } = useAlert();

  function hideAlert() {
    setAlert("");
  }

  if (!text) {
    return <></>;
  }

  return (
    <div className="py-4 px-2">
      <div className="flex bg-red-100 rounded px-2 py-2 text-red-800 mb-2">
        <div className="flex-grow">{text}</div>
        <button onClick={hideAlert}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </button>
      </div>
    </div>
  );
}
