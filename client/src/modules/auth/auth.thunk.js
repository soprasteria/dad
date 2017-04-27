// Imports for fetch API
import { withAuth } from '../auth/auth.wrappers';
import { checkHttpStatus, handleError, parseJSON } from '../utils/promises';

// Auth Actions
import AuthActions from './auth.actions';

// Calls the API to get a token and
// dispatches actions along the way
const loginUser = (auth) => {

  let config = {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: `username=${encodeURIComponent(auth.username)}&password=${encodeURIComponent(auth.password)}`
  };

  return dispatch => {
    // We dispatch requestLogin to kickoff the call to the API
    dispatch(AuthActions.requestLogin());

    return fetch('/auth/login', config)
      .then (checkHttpStatus)
      .then(parseJSON)
      .then((user) =>  {
          // When uer is authorized, add the JWT token in the localstorage for authentication purpose
        localStorage.setItem('id_token', user.id_token);
        dispatch(AuthActions.receiveLogin(user));
      }).catch(error => {
        // When error happens.
        // Dispatch differents actions wether the user is not authorized
        // or if the server encounters any other error
        if (error.response) {
          handleError(error, AuthActions.loginInvalidRequest, dispatch);
        } else {
          dispatch(AuthActions.loginInvalidRequest(error.message));
        }
      });
  };
};

// Logs the user out
const logoutUser = () => {
  return dispatch => {
    dispatch(AuthActions.requestLogout());
    localStorage.removeItem('id_token');
    dispatch(AuthActions.receiveLogout());
  };
};

// Get the profile of the authenticated user
const profile = () => {
  return dispatch => {
    dispatch(AuthActions.requestProfile());

    return fetch('/api/profile', withAuth({ method: 'GET' }))
      .then(checkHttpStatus)
      .then(parseJSON)
      .then(response => {
        dispatch(AuthActions.receiveProfile(response));
      })
      .catch(error => {
        handleError(error, AuthActions.profileError, dispatch);
      });
  };
};

export default {
  loginUser,
  logoutUser,
  profile
};
