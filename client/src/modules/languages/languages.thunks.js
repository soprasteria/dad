import { generateEntitiesThunks } from '../utils/entities';
import LanguagesActions from './languages.actions';

// selectLanguage the language of the app
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
