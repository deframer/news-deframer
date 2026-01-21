import { createRoot } from 'react-dom/client';

import { globalStyles } from '../../shared/theme';
import { Options } from './options';

const init = async () => {
  // Inject global styles
  const style = document.createElement('style');
  style.textContent = globalStyles;
  document.head.appendChild(style);

  // Render the app
  const root = createRoot(document.getElementById('root')!);
  root.render(<Options />);
};

init();
