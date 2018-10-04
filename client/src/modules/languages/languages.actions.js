import { generateEntitiesActions } from '../utils/entities';

import LanguagesConstants from './languages.constants';

// Action when authenticated user successfully get his profile information
const receiveLanguage = (language) => {
  return {
    type: LanguagesConstants.SELECT_LANGUAGE,
    language
  };
};

export default {
  ...generateEntitiesActions('languages'),
  receiveLanguage
};
