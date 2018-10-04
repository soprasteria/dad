import { generateEntitiesReducer } from '../utils/entities';
import LanguagesConstants from './languages.constants';

const languagesReducer = (state, action) => {
  const entitiesState = generateEntitiesReducer(state, action, 'languages');
  switch (action.type) {
  case LanguagesConstants.SELECT_LANGUAGE:
    return {
      ...entitiesState,
      language: action.language
    };
  default:
    return entitiesState;
  }
};

export default languagesReducer;
