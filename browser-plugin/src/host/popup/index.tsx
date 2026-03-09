import '../styles.css';

import { createRoot } from 'react-dom/client';

import { Popup } from './Popup';

const root = createRoot(document.getElementById('root')!);
root.render(<Popup />);
