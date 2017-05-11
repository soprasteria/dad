// import constants
import { REQUEST_INDICATORS, RECEIVE_INDICATORS, INVALID_INDICATORS } from './indicators.constants';

export const Actions = {
  requestSome: (id) => {
    return {
      type: REQUEST_INDICATORS,
      id
    };
  },
  receiveSome: (items) => {
    return {
      type: RECEIVE_INDICATORS,
      items,
      receivedAt: Date.now()
    };
  },
  invalidRequestEntity: (items) => (error) => {
    return {
      type: INVALID_INDICATORS,
      title: 'Cannot fetch indicators for this project',
      message: error,
      level: 'error',
      items
    };
  },
};


export default {
  Actions
};
