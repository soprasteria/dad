// import constants
import EntitiesConstants from './entities.constants';
import { generateEntitiesActions } from '../utils/entities';


// Change filter
const changeFilter = (filterValue) => {
  return {
    type: EntitiesConstants.CHANGE_FILTER,
    filterValue
  };
};

export default {
  ...generateEntitiesActions('entities'),
  changeFilter
};
