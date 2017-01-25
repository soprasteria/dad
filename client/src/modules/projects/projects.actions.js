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

// Add URL
const addUrl = (id, url) => {
  return {
    type: ProjectsConstants.ADD_URL,
    id,
    url
  };
};

// Edit URL
const editUrl = (id, index, url) => {
  return {
    type: ProjectsConstants.EDIT_URL,
    id,
    index,
    url
  };
};

// Edit URL
const removeUrl = (id, index) => {
  return {
    type: ProjectsConstants.REMOVE_URL,
    id,
    index
  };
};

export default {
  ...generateEntitiesActions('projects'),
  changeFilter,
  addUrl,
  editUrl,
  removeUrl
};
