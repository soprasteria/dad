import React from 'react';
import PropTypes from 'prop-types';
import { Provider } from 'react-redux';
import { IndexRoute, Route, Router } from 'react-router';

import App from './app.layout';
import Home from './home.page';
import ProjectsPage from '../projects/projects.page';
import ProjectPage from '../projects/project/project.page';
import UsersPage from '../users/users.page';
import UserPage from '../users/user/user.page';

import { requireAuthorization } from '../auth/auth.isAuthorized';
import { AUTH_ADMIN_ROLE, AUTH_RI_ROLE } from '../../modules/auth/auth.constants';

const RootComponent = ({
  store,
  history
}) => (
  <Provider store={store}>
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
  </Provider>
);

RootComponent.propTypes = {
  store: PropTypes.object.isRequired,
  history: PropTypes.object.isRequired
};

export default RootComponent;
