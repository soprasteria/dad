import { generateEntitiesConstants } from '../utils/entities';

export default {
  ...generateEntitiesConstants('entities'),
  CHANGE_FILTER: 'CHANGE_FILTER_ENTITIES'
};
