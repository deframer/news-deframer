export const Spinner = () => {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
      }}
    >
      <style>
        {`
          @keyframes spin {
            from { transform: rotate(0deg); }
            to { transform: rotate(360deg); }
          }
        `}
      </style>
      <img
        src="chrome-extension://__MSG_@@extension_id__/assets/spinner.svg"
        alt="Loading..."
        style={{
          width: '64px',
          height: '64px',
          animation: 'spin 1.5s linear infinite',
        }}
      />
    </div>
  );
};
