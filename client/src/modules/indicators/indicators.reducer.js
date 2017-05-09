// import constants
import IndicatorsConstants from './indicators.constants';
import { generateEntitiesReducer } from '../utils/entities';

const indicatorsReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'indicators');
  switch (action.type) {
  case IndicatorsConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  default:
    return entitiesState;
  }
};

export default indicatorsReducer;
