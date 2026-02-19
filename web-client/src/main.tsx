import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'urql';
import App from './App';
import { client } from './graphql/client';
import { AuthProvider } from './context/AuthContext';
import './styles/global.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Provider value={client}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </Provider>
  </React.StrictMode>,
);
