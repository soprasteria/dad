// import constants
import OrganizationsConstants from './organizations.constants';
import { generateEntitiesReducer } from '../utils/entities';

const organizationsReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'organizations');
  switch (action.type) {
  case OrganizationsConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  default:
    return entitiesState;
  }
};

export default organizationsReducer;
