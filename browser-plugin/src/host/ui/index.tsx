import './styles.css';

import { createRoot } from 'react-dom/client';

import { Options } from './options';

const init = async () => {
  // Render the app
  const root = createRoot(document.getElementById('root')!);
  root.render(<Options />);
};

init();
