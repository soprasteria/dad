// import constants
import { INVALID_INDICATORS, REQUEST_INDICATORS, RECEIVE_INDICATORS } from './indicators.constants';
import { initialState } from '../utils/entities';

const indicatorsReducer = (state = initialState, action) => {
  switch (action.type) {
  case INVALID_INDICATORS:
    return {
      ...state,
      ...initialState,
      items: { ...state.items }
    };
  case REQUEST_INDICATORS:
    return {
      ...state,
      isFetching: true,
      didInvalidate: false
    };
  case RECEIVE_INDICATORS:
    let items = {};
    action.items.forEach((item) => items[item.id] = { ...state.items[item.id], ...item });
    return {
      ...state,
      isFetching: false,
      didInvalidate: false,
      items,
      lastUpdated: action.receivedAt
    };
  default:
    return state;  
  };
};


export default indicatorsReducer;
