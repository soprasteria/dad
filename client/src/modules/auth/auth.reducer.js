// import constants
import AuthConstants from './auth.constants';
import UsersConstants from '../users/users.constants';
import { LOCATION_CHANGE } from 'react-router-redux';


const initialState = {
  isFetching: false,
  isAuthenticated: localStorage.getItem('id_token') ? true : false,
  user: {}
};

// The auth reducer. The starting state sets authentication
// based on a token being in local storage. In a real app,
// we would also want a util to check if the token is expired.
const authReducer = (state = initialState, action) => {
  switch (action.type) {
  case LOCATION_CHANGE:
    return { ...state, errorMessage: '' };
  case AuthConstants.LOGIN_REQUEST:
    return {
      ...state,
      isFetching: true,
      isAuthenticated: false,
      user: {}
    };
  case AuthConstants.LOGIN_SUCCESS:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: true,
      errorMessage: '',
      user: action.user
    };
  case AuthConstants.LOGIN_INVALID_REQUEST:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: false,
      user: {}
    };
  case AuthConstants.LOGIN_NOT_AUTHORIZED:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: false,
      errorMessage: action.error,
      user: {}
    };
  case AuthConstants.LOGOUT_SUCCESS:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: false,
      user: {}
    };
  case AuthConstants.PROFILE_REQUEST:
    return {
      ...state,
      isFetching: true,
    };
  case AuthConstants.PROFILE_SUCCESS:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: true,
      errorMessage: '',
      user: action.user
    };
  case AuthConstants.PROFILE_FAILURE:
    return {
      ...state,
      isFetching: false,
      isAuthenticated: false,
      errorMessage: action.message,
      user: {}
    };
  case UsersConstants.REQUEST_SAVE_USER:
    return { ...state, ...authenticatedUserIsFetching(state, action) };
  case UsersConstants.USER_SAVED:
    return { ...state, ...changeUserIfNeeded(state, action) };
  case UsersConstants.INVALID_SAVE_USER:
    return { ...state, ...authenticatedUserFetchingError(state, action) };
  case UsersConstants.REQUEST_DELETE_USER:
    if (action.id === state.user.id) {
      return {
        ...state,
        user: { ...state.user, isDeleting: true },
      };
    } else {
      return state;
    }
  case UsersConstants.USER_DELETED:
    if (action.id === state.user.id) {
      return {
        isAuthenticated: false,
        user: {},
        isFetching: false
      };
    } else {
      return state;
    }
  case UsersConstants.INVALID_DELETE_USER:
    if (action.entity.id === state.user.id) {
      return {
        ...state,
        user: { ...state.user, isDeleting: false },
      };
    } else {
      return state;
    }
  default:
    return state;
  }
};

const authenticatedUserIsFetching = (state, action) => {
  if (action.entity.id === state.user.id) {
    return {
      user: { ...action.entity, isFetching: true, errorMessage: '' }
    };
  } else {
    return {};
  }
};

const authenticatedUserFetchingError = (state, action) => {
  if (action.entity.id === state.user.id) {
    return {
      user: { ...action.entity, isFetching: false, errorMessage: action.error }
    };
  } else {
    return {};
  }
};

const changeUserIfNeeded = (state, action) => {
  if (action.entity.id === state.user.id) {
    return {
      user: { ...action.entity, isFetching: false, errorMessage: '' }
    };
  } else {
    return {};
  }
};

export default authReducer;
