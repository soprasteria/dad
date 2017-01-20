// import constants
import EntitiesConstants from './entities.constants';
import { generateEntitiesReducer } from '../utils/entities';

const entitiesReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'entities');
  switch (action.type) {
  case EntitiesConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  default:
    return entitiesState;
  }
};

export default entitiesReducer;
