import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

const textarea = document.querySelector("#app-host") as HTMLInputElement
let host: string

if (!textarea) {
    host = "http://localhost:8080"
} else {
    host = textarea.value
}

root.render(
    <React.StrictMode>
        <App host={host}/>
    </React.StrictMode>
);

