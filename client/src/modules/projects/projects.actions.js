// import constants
import ProjectsConstants from './projects.constants';
import { generateEntitiesActions } from '../utils/entities';


// Change filter
const changeFilter = (filterValue) => {
  return {
    type: ProjectsConstants.CHANGE_FILTER,
    filterValue
  };
};

export default {
  ...generateEntitiesActions('projects'),
  changeFilter
};
