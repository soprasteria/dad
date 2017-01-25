// import constants
import ProjectsConstants from './projects.constants';
import { generateEntitiesReducer } from '../utils/entities';

const projectsReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'projects');
  switch (action.type) {
  case ProjectsConstants.CHANGE_FILTER:
    return { ...entitiesState, filterValue: action.filterValue };
  case ProjectsConstants.ADD_URL:
    const projectToAdd = { ...entitiesState.items[action.id] };
    projectToAdd.urls = projectToAdd.urls || [];
    projectToAdd.urls.push(action.url);
    projectToAdd.isEditing = true;
    return {
      ...entitiesState,
      items: { ...entitiesState.items, [action.id]:projectToAdd }
    };
  case ProjectsConstants.EDIT_URL:
    const projectToEdit = { ...entitiesState.items[action.id] };
    const urlToEdit = projectToEdit.urls[action.index];
    if (urlToEdit) {
      projectToEdit.urls[action.index] = action.url;
    }
    projectToEdit.isEditing = true;
    return {
      ...entitiesState,
      items: { ...entitiesState.items, [action.id]:projectToEdit }
    };
  case ProjectsConstants.REMOVE_URL:
    const projectToRemove = { ...entitiesState.items[action.id] };
    projectToRemove.urls.splice(action.index, 1);
    projectToRemove.isEditing = true;
    return {
      ...entitiesState,
      items: { ...entitiesState.items, [action.id]:projectToRemove }
    };
  default:
    const projects = { ...entitiesState, items: { ...entitiesState.items } };
    Object.values(projects.items).forEach(item => item = { ...item, isEditing: false });
    return projects;
  }
};

export default projectsReducer;
