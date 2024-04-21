import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import {AuthAPI} from "./api/AuthAPI";
import {ExpressionsAPI} from "./api/ExpressionsAPI";
import {OperatorsAPI} from "./api/OperatorsAPI";
import {WorkerAPI} from "./api/WorkersAPI";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const authAPI = new AuthAPI();
const expressionsAPI = new ExpressionsAPI();
const operatorsAPI = new OperatorsAPI();
const workersAPI = new WorkerAPI();

root.render(
  <React.StrictMode>
    <App
        authAPI={authAPI}
        expressionsAPI={expressionsAPI}
        operatorsAPI={operatorsAPI}
        workersAPI={workersAPI}
    />
  </React.StrictMode>
);
