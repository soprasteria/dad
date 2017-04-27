import { withAuth } from '../auth/auth.wrappers';
import { checkHttpStatus, handleError } from '../utils/promises';
import { download } from '../utils/files';

// Export actions
import ExportActions from './export.actions';

import moment from 'moment';

// Calls the API to export data
const exportAll = () => {

  return dispatch => {
    // We dispatch requestLogin to kickoff the call to the API
    dispatch(ExportActions.requestExportAll());

    return fetch('/api/export', withAuth({ method: 'GET' }))
      .then(checkHttpStatus)
      .then(response => {
        response.blob().then(blob => {
          const currentDate = moment(new Date()).format('YYYY-MM-DD-HH-mm-ss');
          download(blob, `deployment-plan-${currentDate}.xlsx`);
          dispatch(ExportActions.receiveExportAll());
        });
      })
      .catch(error => {
        handleError(error, ExportActions.invalidRequestExportAll, dispatch);
      });
  };
};

export default {
  exportAll
};
