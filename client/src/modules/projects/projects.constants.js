import { generateEntitiesConstants } from '../utils/entities';

export default {
  ...generateEntitiesConstants('projects'),
  CHANGE_FILTER: 'CHANGE_FILTER_PROJECTS'
};
