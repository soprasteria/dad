import { generateEntitiesConstants } from '../utils/entities';

export default {
  ...generateEntitiesConstants('projects'),
  CHANGE_FILTER: 'CHANGE_FILTER_PROJECTS',
  ADD_URL: 'ADD_URL',
  EDIT_URL: 'EDIT_URL',
  REMOVE_URL: 'REMOVE_URL'
};
