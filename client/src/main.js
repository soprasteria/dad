import 'babel-polyfill';
import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { IndexRoute, Route, Router, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

// Store
import { store } from './store';

// Components
import App from './components/app/app.layout';
import Home from './components/app/home.page';
import ProjectsPage from './components/projects/projects.page';
import ProjectPage from './components/projects/project/project.page';
import UsersPage from './components/users/users.page';
import UserPage from './components/users/user/user.page';
import { requireAuthorization } from './components/auth/auth.isAuthorized';

// Constants
import { AUTH_ADMIN_ROLE, AUTH_RI_ROLE } from './modules/auth/auth.constants';

// Create an enhanced history that syncs navigation events with the store
const history = syncHistoryWithStore(browserHistory, store);

ReactDOM.render(
  <Provider store={store}>
    {/* Tell the Router to use our enhanced history */}
    <Router history={history}>
      <Route path='/' component={App}>
        <IndexRoute component={Home} />
        <Route path='login' component={Home} />
        <Route path='projects'>
          <IndexRoute component={requireAuthorization(ProjectsPage)} />
          <Route path='new' component={requireAuthorization(ProjectPage, [AUTH_ADMIN_ROLE, AUTH_RI_ROLE])} />
          <Route path=':id' component={requireAuthorization(ProjectPage)} />
        </Route>
        <Route path='users'>
          <IndexRoute component={requireAuthorization(UsersPage)} />
          <Route path=':id' component={requireAuthorization(UserPage)} />
        </Route>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('root')
);
