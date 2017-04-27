// import constants
import ProjectsConstants from './projects.constants';
import ModalConstants from '../modal/modal.constants';
import { generateEntitiesReducer } from '../utils/entities';

const projectsReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'projects');
  switch (action.type) {
  case ProjectsConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  case ModalConstants.OPEN_MODAL:
    let id = entitiesState.selected.id;
    let item = entitiesState.items[id];
    if (!item) {
      return entitiesState;
    }
    item.isEditing = true;
    return {
      ...entitiesState,
      items: { ...entitiesState.items, [id]: item }
    };

  default:
    const projects = { ...entitiesState, items: { ...entitiesState.items } };
    Object.values(projects.items).forEach(item => item = { ...item, isEditing: false });
    return projects;
  }
};

export default projectsReducer;
