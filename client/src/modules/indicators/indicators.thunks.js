import { withAuth } from '../auth/auth.wrappers';
import { checkHttpStatus, handleError, parseJSON } from '../utils/promises';
import Actions from './indicators.actions';

export const fetchIndicators = (id) => {
  return function (dispatch) {
    dispatch(Actions.requestSome(id));
    return fetch(`/api/projects/${id}/indicators`, withAuth({ method: 'GET' }))
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => {
        dispatch(Actions.receiveSome(response));
      })
      .catch((error) => {
        handleError(error, Actions.invalidRequestEntity({ id }), dispatch);
      });
  };
};
