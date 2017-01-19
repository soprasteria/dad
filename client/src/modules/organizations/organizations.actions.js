// import constants
import OrganizationsConstants from './organizations.constants';
import { generateEntitiesActions } from '../utils/entities';


// Change filter
const changeFilter = (filterValue) => {
  return {
    type: OrganizationsConstants.CHANGE_FILTER,
    filterValue
  };
};

export default {
  ...generateEntitiesActions('organizations'),
  changeFilter
};
