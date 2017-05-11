// import constants
import IndicatorsConstants from './indicators.constants';

export default {
  requestSome: (id) => {
    return {
      type: IndicatorsConstants.REQUEST_INDICATORS,
      id
    };
  },
  receiveSome: (items) => {
    return {
      type: IndicatorsConstants.RECEIVE_INDICATORS,
      items,
      receivedAt: Date.now()
    };
  },
  invalidRequestEntity: (items) => (error) => {
    return {
      type: IndicatorsConstants.INVALID_INDICATORS,
      title: 'Cannot fetch indicators for this project',
      message: error,
      level: 'error',
      items
    };
  },
};
