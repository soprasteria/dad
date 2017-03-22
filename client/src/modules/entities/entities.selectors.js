import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';

export const getFilteredEntities = (entities, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(entities);
  } else {
    return Object.values(entities).filter(entity => {
      return containsWithoutAccents(JSON.stringify(Object.values(entity)), filterValue);
    });
  }
};

export const getByType = (entities, type) => {
  return entities.filter(entity => entity.type === type);
};

export const getEntitiesAsOptions = ({ entities, withType }) => {
  return [{ value: '', text:'None' }].concat(entities.map(entity => {
    let type = entity.type;
    type = type.charAt(0).toUpperCase() + type.substring(1);
    const value = withType ? `${type}: ${entity.name}` : entity.name;
    return { value: entity.id, text: value } ;
  }).sort(sortby('text')));
};
