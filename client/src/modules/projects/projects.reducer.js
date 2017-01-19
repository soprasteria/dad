// import constants
import ProjectsConstants from './projects.constants';
import { generateEntitiesReducer } from '../utils/entities';

const projectsReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'projects');
  switch (action.type) {
  case ProjectsConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  default:
    return entitiesState;
  }
};

export default projectsReducer;
