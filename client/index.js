import 'babel-polyfill';

import React from 'react';
import ReactDOM from 'react-dom';

import App from './App';

console.log(process.env.NODE_ENV);
ReactDOM.render(<App />, document.getElementById('app'));
