import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

// Store
import { store } from './store';

// Components
import App from './components/app/app.layout';
import Home from './components/app/home.page';
import UsersPage from './components/users/users.page';
import { requireAuthorization } from './components/auth/auth.isAuthorized';

// Create an enhanced history that syncs navigation events with the store
const history = syncHistoryWithStore(browserHistory, store);

ReactDOM.render(
  <Provider store={store}>
    {/* Tell the Router to use our enhanced history */}
    <Router history={history}>
      <Route path='/' component={App}>
        <IndexRoute component={Home} />
        <Route path='login' component={Home} />
        <Route path='users'>
          <IndexRoute component={requireAuthorization(UsersPage)} />
        </Route>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('root')
);
