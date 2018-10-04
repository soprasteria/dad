import { generateEntitiesThunks } from '../utils/entities';
import LanguagesActions from './languages.actions';

// Logs the user out
const selectLanguage = (language) => {
  return (dispatch) => {
    dispatch(LanguagesActions.receiveLanguage(language));
    localStorage.setItem('language', language);
  };
};

export default {
  ...generateEntitiesThunks('languages'),
  selectLanguage
};
