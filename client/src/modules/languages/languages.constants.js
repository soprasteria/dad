import { generateEntitiesConstants } from '../utils/entities';

export default {
  ...generateEntitiesConstants('languages'),
  SELECT_LANGUAGE: 'SELECT_LANGUAGE',
  DEFAULT_LANGUAGE: 'fr'
};
