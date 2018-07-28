import 'babel-polyfill';

import React from 'react';
import ReactDOM from 'react-dom';

import App from './App';

console.log(process.env.SERVER_URL);
ReactDOM.render(<App />, document.getElementById('app'));
