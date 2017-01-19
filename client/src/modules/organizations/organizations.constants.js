import { generateEntitiesConstants } from '../utils/entities';

export default {
  ...generateEntitiesConstants('organizations'),
  CHANGE_FILTER: 'CHANGE_FILTER_ORGANIZATIONS'
};
