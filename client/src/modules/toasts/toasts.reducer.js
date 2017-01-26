//JS dependancies
import UUID from 'uuid-js';
import MD5 from 'md5';

//Actions
import { LOCATION_CHANGE } from 'react-router-redux';
import EntitiesConstants from '../entities/entities.constants';
import ProjectsConstants from '../projects/projects.constants';
import UsersConstants from '../users/users.constants';
import ServicesConstants from '../services/services.constants';
import AuthConstants from '../auth/auth.constants';
import ExportConstants from '../export/export.constants';
import ToastsConstants from './toasts.constants';

const initialState = {};

const toastsReducer = (state = initialState, action) => {
  switch (action.type) {
  case EntitiesConstants.INVALID_REQUEST_ENTITIES:
  case EntitiesConstants.INVALID_REQUEST_ENTITY:
  case EntitiesConstants.INVALID_SAVE_ENTITY:
  case EntitiesConstants.INVALID_DELETE_ENTITY:
  case ProjectsConstants.INVALID_REQUEST_PROJECTS:
  case ProjectsConstants.INVALID_REQUEST_PROJECT:
  case ProjectsConstants.INVALID_SAVE_PROJECT:
  case ProjectsConstants.INVALID_DELETE_PROJECT:
  case ServicesConstants.INVALID_REQUEST_SERVICES:
  case ServicesConstants.INVALID_REQUEST_SERVICE:
  case ServicesConstants.INVALID_SAVE_SERVICE:
  case ServicesConstants.INVALID_DELETE_SERVICE:
  case UsersConstants.INVALID_REQUEST_USERS:
  case UsersConstants.INVALID_REQUEST_USER:
  case UsersConstants.INVALID_SAVE_USER:
  case UsersConstants.INVALID_DELETE_USER:
  case AuthConstants.LOGIN_INVALID_REQUEST:
  case ExportConstants.EXPORT_ALL_INVALID_REQUEST:
    return { ...state, ...createGenericToast(action) };
  case ToastsConstants.COMFIRM_DELETION:
    return { ...state, ...createConfirmDelToast(action) };
  case ToastsConstants.CLOSE_NOTIFICATION:
    let resState = { ...state };
    delete resState[action.id];
    return resState;
  case LOCATION_CHANGE:
    return { ...initialState };
  default:
    return state;
  }
};

const createGenericToast = (action) => {
  let res = {};
  const uuid = UUID.create(4).hex;
  res[uuid] = {
    title: action.title,
    message: action.message,
    level: action.level,
    autoDismiss: 10,
    position: 'br',
    uid: uuid
  };
  return res;
};

const createConfirmDelToast = (action) => {
  let res = {};
  const id = MD5(action.title);
  res[id] = {
    title: 'Confirm Suppression',
    message: 'Remove ' + action.title + ' ?',
    autoDismiss: 0,
    level: 'error',
    position: 'br',
    uid: id,
    action: {
      label: 'Remove',
      callback: action.callback
    }
  };
  return res;
};

export default toastsReducer;
