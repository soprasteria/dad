import 'babel-polyfill';
import 'isomorphic-fetch';
import React from 'react';
import ReactDOM from 'react-dom';

import { AppContainer } from 'react-hot-loader';
import RootComponent from './components/app/root';

import { store } from './store';
import { browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';
const history = syncHistoryWithStore(browserHistory, store);

const render = (Root) => {
  ReactDOM.render(
    <AppContainer>{Root}</AppContainer>,
    document.getElementById('root')
  );
};

render(<RootComponent history={history} store={store} />);

if (module.hot) {
  module.hot.accept('./components/app/root', () => {
    const RootComponent = require('./components/app/root');
    render(<RootComponent history={history} store={store} />);
  });
}
